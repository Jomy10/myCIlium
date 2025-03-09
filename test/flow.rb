require 'net/http'
require 'json'

puts "tests"

$tests << [
  "flow",
  proc { |t|

    puts "== Request Build =="

    uri = URI("#{HOST}/request-build")
    http = Net::HTTP::new(uri.host, uri.port)
    http.set_debug_output($stdout)
    req = Net::HTTP::Post.new(uri.request_uri)
    req["Authorization"] = "Bearer #{CREATE_TOKEN}"
    req.body = '{"repo": "https://github.com/user/repo.git", "platforms": ["macos", "linux"], "rev": "some-commit"}'

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

    puts "== Start building =="

    uri = URI("#{HOST}/request-start")
    req = Net::HTTP::Post.new(uri.request_uri)
    req["Authorization"] = "Bearer #{MACOS_TOKEN}"
    req.body = '{"requestId": 1}'

    resp = http.request(req)
    if resp.code.to_i != 200
      t.fail("Found code #{resp.code}")
      next
    end

    puts "== Try building same id =="

    puts "== Finish Build =="

    puts "== Try finishing same id =="

  }
]
