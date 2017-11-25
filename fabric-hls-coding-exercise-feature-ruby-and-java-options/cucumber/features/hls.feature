Feature: HLS Manifests
  As Brightcove
  I Want to serve HLS manifests
  So That happy

  @task2
  Scenario: Healthcheck
    When I call GET on "/healthcheck"
    Then the result should be a "200"

  @task3
  Scenario: Metadata API
    Given the metadata file "simple.json"
    When I call GET on "/metadata/simple"
    Then the result should be a "200"
    And the response should be JSON equivalent to "simple.json"
    And the content type should be "application/json"

  @task4
  Scenario: Duration API
    Given the metadata file "simple.json"
    When I call GET on "/duration/simple"
    Then the result should be a "200"
    And the response should match "simple-duration.txt"
    And the content type should be "text/plain"

  @task5
  Scenario: Simple playback
    Given the metadata file "simple.json"
    When I call GET on "/manifest/simple/2s.m3u8"
    Then the result should be a "200"
    And the response should match "expected-manifest-2s-seg.m3u8"
    And the content type should be "application/x-mpegURL"

  @task6
  Scenario: 404s on non-existent ID
    When I call GET on "/manifest/idonotexist/2s.m3u8"
    Then the result should be a "404"

  @task7
  Scenario: Variable segment duration
    Given the metadata file "simple.json"
    When I call GET on "/manifest/simple/10s.m3u8"
    Then the result should be a "200"
    And the response should match "expected-manifest-10s-seg.m3u8"
    And the content type should be "application/x-mpegURL"
