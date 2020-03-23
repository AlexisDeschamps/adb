Design of sendy sync
====================

We want to have a 1-way sync between emails in the ADB and emails in
Sendy.

Cases:

 - when someone is added to the database, add their email to sendy
 - if someone's email is changed, unsubscribe their previous email (if it's the only not already in sendy)


Design:

Use the table supporters_sendy_sync with the following columns:

 - supporter_id
 - sendy_list_id
 - email
 - sync_status: an enum, either:

1 - exists in sendy (either because it was added successfully or b/c it already existed)
2 - non-retryable error (i.e. there's some issue with adding it to sendy, and we're not going to retry)

 - sync_timestamp: timestamp for when they were added

TODO: Unsubscribe & resubscribe new email if someone's email is updated.
