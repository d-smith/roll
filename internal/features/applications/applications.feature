@apptests
Feature: Application Tests

  Scenario: Application Registration
    Given a developer registered with the portal
    And they have a new application they wish to register
    Then the application should be successfully registered

  Scenario: Retrieve Application Details
    Given a registered application
    Then the details associated with the application can be retrieved

  Scenario: Duplicate Application Registration
    Given an application has already been registered
    And a developer attempts to register an application with the same name
    Then an error is returned with status code StatusConflict
    And the error message indicates a duplicate registration was attempted

  Scenario: Application Registration Update
    Given a registered application to update
    And there are updates to make to the application defnition
    Then the application can be updated
    And the updates are reflected when retrieving the application definition anew