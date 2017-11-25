Given(/^the metadata file "([^"]*)"$/) do |file|
end

When(/^I call GET on "(.*)"$/) do |path|
  @response = HTTParty.get("http://localhost:1337#{path}")
end

Then(/^the result should be a "(\d+)"$/) do |stat_code|
  # We _really_ shouldn't need code.to_i, but when @response has been pulled out of an exception, @response.code
  # becomes a string. Sigh.
  puts @response.body if @response.code.to_i != stat_code.to_i
  @response.code.to_i.should eq stat_code.to_i
end

And(/^the content type should be "(.*)"$/) do |content_type|
  @response.headers['content-type'].should eq(content_type)
end

And(/^the response should match "(.*)"$/) do |fixture|
  fixture_contents = File.read(File.join(CUCUMBER_BASE, 'fixtures', fixture))
  @response.body.should eq(fixture_contents)
end

And(/^the response should be JSON equivalent to "(.*)"$/) do |fixture|
  fixture_contents = File.read(File.join(CUCUMBER_BASE, 'fixtures', fixture))
  @response.body.should be_json_eql(fixture_contents)
end
