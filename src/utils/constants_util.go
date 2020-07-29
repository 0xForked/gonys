package utils

/*
 * Gonys - A Notification Service for SMS
 *
 * Constants Utilities
 *
 * @author A. A. Sumitro <hello@aasumitro.id>
 * https://aasumitro.id
 */

// Delivery report announced pending deliver.
const DeliveryPending = "DELIVERY_PENDING" //  queueing status

// Some other error happened during sending.
const SendingError = "SENDING_ERROR" // retrying status

// Message has been sent, waiting for delivery report.
const SendingOk = "SENDING_OK" // complete status
