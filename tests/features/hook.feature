Feature: Hooks
    In order to handle payment Hooks
    As a gateway payment
    I want to be able to send a payment hook

    Scenario: Send a payment hook for a successful payment
        Given I have a payment
        When I send a payment hook for a "successful" payment
        Then the payment state should be "Approved"

    Scenario: Send a payment hook for a failed payment
        Given I have a payment
        When I send a payment hook for a "failed" payment
        Then the payment state should be "Rejected"