require 'net/ping'
require 'sqlite3'

$tests = []

HOST = "http://localhost:8080"
CREATE_TOKEN = "CREATE"
MACOS_TOKEN = "MACOS"

class Tester
  def initialize name
    @name = name
    @failed = []
    @succeeded = []
  end

  def assert cond, msg = nil
    if cond
      @succeeded << [@name, msg, caller_locations.first]
    else
      @failed << [@name, msg, caller_locations.first]
    end
  end

  def fail(msg = nil)
    @failed << [@name, msg, caller_locations]
  end

  def str
    str = "#{@name}: #{@failed.size} failed, #{@succeeded.size} succeeded"
    for failed in @failed
      str << "\nfailed: #{failed}"
    end
    return str
  end
end

pid = spawn("go run .")

# Wait for server to start
while true
  begin
    uri = URI(HOST)
    http = Net::HTTP::new(uri.host, uri.port)
    req = Net::HTTP::Get.new(uri.request_uri)
    http.request(req)
    break
  rescue
    next
  end
end

puts "Server started and accepting connections"

begin
  db = SQLite3::Database.new "data.db"

  # Create debug token
  sql = <<-SQL
  select count(*)
  from Token
  where token = ?
  SQL
  rows = db.execute sql, CREATE_TOKEN
  if rows.first.first == 0
    sql = <<-SQL
    insert into Token (token, right)
    values (?, 'create')
    SQL
    p db.execute sql, CREATE_TOKEN
  end

  sql = <<-SQL
  select count(*)
  from Token
  where token = ?
  SQL
  rows = db.execute sql, MACOS_TOKEN
  if rows.first.first == 0
    sql = <<-SQL
    insert into Token (token, platform, right)
    values (?, 'macos', 'platform')
    SQL
    p db.execute sql, MACOS_TOKEN
  end
rescue => e
  p e
  puts "closing connection"
  Process.kill("INT", pid)
  return
end

begin
  Dir["test/**/*.rb"].each do |file|
    require_relative file
  end

  res = []
  for test in $tests
    border = "#" * (test[0].size + 4)
    puts border
    puts "# #{test[0]} #"
    puts border

    tester = Tester.new test[0]
    begin
      test[1].call tester
      res << tester
    rescue => e
      tester.fail(e.to_s)
      res << tester
    end
  end
ensure
  puts "closing connection"
  Process.kill("INT", pid)
end

for res in res
  puts res.str
end
