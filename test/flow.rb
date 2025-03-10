require 'net/http'
require 'json'

puts "tests"

# This test follows the complete flow of
# 1. Creating tasks
# 2. Polling tasks
# 3. Starting a task
# 4. Finishing a task
$tests << [
  "flow",
  proc { |t|
    repo = "https://github.com/user/repo.git"

    puts "== Request Build =="

    uri = URI("#{HOST}/request-build")
    http = Net::HTTP::new(uri.host, uri.port)
    http.set_debug_output($stdout)
    req = Net::HTTP::Post.new(uri.request_uri)
    req["Authorization"] = "Bearer #{CREATE_TOKEN}"
    req.body = '{"repo": "' + repo + '", "platforms": ["macos", "linux"], "rev": "some-commit"}'

    resp = http.request(req)
    if resp.code.to_i != 201
      t.fail("Found code #{resp.code}")
      next
    end
    resp = JSON::parse(resp.body)
    ids = resp["ids"]

    t.assert(ids.size == 2, "ids size #{ids.size} != 2")

    macosId = ids.find do |req|
      req["platform"] == "macos"
    end
    macosId = macosId["requestId"]

    puts "macosId: #{macosId}"
    t.assert(macosId != nil)

    puts "== Poll tasks =="
    uri = URI("#{HOST}/requests")
    req = Net::HTTP::Post.new(uri.request_uri)
    req.body = "{\"platform\": \"macos\", \"status\": \"requested\"}"

    resp = http.request(req)
    if resp.code.to_i != 200
      t.fail("Found code #{resp.code}")
      next
    end
    resp = JSON::parse(resp.body)

    found = false
    for req in resp["requests"]
      t.assert(req["status"] == 1 && req["platform"] == "macos")
      if req["id"] == macosId && req["repo"] == repo && req["revision"] == "some-commit"
        found = true
      end
    end
    t.assert(found)

    puts "== Start building =="

    uri = URI("#{HOST}/request-start")
    req = Net::HTTP::Post.new(uri.request_uri)
    req["Authorization"] = "Bearer #{MACOS_TOKEN}"
    req.body = "{\"requestId\": #{macosId}}"

    resp = http.request(req)
    if resp.code.to_i != 200
      t.fail("Found code #{resp.code}")
      next
    end

    puts "== Check if status updated =="

    uri = URI("#{HOST}/requests")
    req = Net::HTTP::Post.new(uri.request_uri)
    req.body = "{\"platform\": \"macos\", \"status\": \"started\"}"

    resp = http.request(req)
    if resp.code.to_i != 200
      t.fail("Found code #{resp.code}")
      next
    end
    resp = JSON::parse(resp.body)

    found = false
    for req in resp["requests"]
      t.assert(req["status"] == 2 && req["platform"] == "macos")
      if req["id"] == macosId && req["repo"] == repo && req["revision"] == "some-commit"
        found = true
      end
    end
    t.assert(found)

    puts "== Try building same id =="

    uri = URI("#{HOST}/request-start")
    req = Net::HTTP::Post.new(uri.request_uri)
    req["Authorization"] = "Bearer #{MACOS_TOKEN}"
    req.body = "{\"requestId\": #{macosId}}"

    resp = http.request(req)
    if resp.code.to_i != 409
      t.fail("Found code #{resp.code}")
      next
    end

    puts "== Finish Build =="

    uri = URI("#{HOST}/request-finish")
    req = Net::HTTP::Post.new(uri.request_uri)
    req["Authorization"] = "Bearer #{MACOS_TOKEN}"
    req.body = "{\"requestId\": #{macosId}}"

    resp = http.request(req)
    if resp.code.to_i != 200
      t.fail("Found code #{resp.code}")
      next
    end

    puts "== Check if status updated to finished =="

    uri = URI("#{HOST}/requests")
    req = Net::HTTP::Post.new(uri.request_uri)
    req.body = "{\"platform\": \"macos\", \"status\": \"finished\"}"

    resp = http.request(req)
    if resp.code.to_i != 200
      t.fail("Found code #{resp.code}")
      next
    end
    resp = JSON::parse(resp.body)

    found = false
    for req in resp["requests"]
      t.assert(req["status"] == 3 && req["platform"] == "macos", "query is wrong")
      if req["id"] == macosId && req["repo"] == repo && req["revision"] == "some-commit"
        found = true
      end
    end
    t.assert(found)

    puts "== Try finishing same id =="

    uri = URI("#{HOST}/request-finish")
    req = Net::HTTP::Post.new(uri.request_uri)
    req["Authorization"] = "Bearer #{MACOS_TOKEN}"
    req.body = "{\"requestId\": #{macosId}}"

    resp = http.request(req)
    if resp.code.to_i != 409
      t.fail("Found code #{resp.code}")
      next
    end

    puts "End flow"
  }
]
