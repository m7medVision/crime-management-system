# TODO List for BE Rihal Codestacker Challenge 2025

## User Management

- [ ] Develop API for admin to add users.
- [ ] Develop API for admin to update users.
- [ ] Develop API for admin to delete users.
- [ ] Admin can assign roles and clearance levels to users.

## Case Management

### Public APIs
- [ ] Develop API to submit a crime report and return the report ID.
- [ ] Develop API to return the status of a case given its report ID.

### Admin APIs
- [ ] Create API to add new cases.
- [ ] Ensure new cases are linked to the user who created them.
- [ ] Develop API to update existing cases.
- [ ] Develop API to return a list of all cases.
  - [ ] Include: Case Number, Case Name, Description, Area/City, Created By/At, Case Type, Authorization Level.
  - [ ] Ensure descriptions are truncated to 100 characters if they exceed the limit.
  - [ ] Allow searching for cases by name or description.

### Case Details
- [ ] Develop API to return detailed information for a specific case by ID.
  - [ ] Include: Case Number, Case Name, Description, Area/City, Created By/At, Case Type, Case Level, Authorization Level, Reported By, Number of Assignees, Number of Evidences, Number of Suspects, Number of Victims, Number of Witnesses.

### Additional Case APIs
- [ ] Develop API to return all assignees of a case by ID.
- [ ] Develop API to return all evidence of a case by ID.
- [ ] Develop API to return all suspects of a case by ID.
- [ ] Develop API to return all victims of a case by ID.
- [ ] Develop API to return all witnesses of a case by ID.

## Evidence Management

- [ ] Develop API to record text or image evidence related to a case.
  - [ ] Include optional remarks.
  - [ ] Validate images to ensure they are actual images.
- [ ] Develop API to retrieve an evidence entry by ID.
  - [ ] Return size of the image if the evidence is an image.
- [ ] Develop API to retrieve the evidence image by ID.
  - [ ] Handle cases where evidence is not an image.
- [ ] Develop API to update an evidence entry.
  - [ ] Ensure the type of evidence cannot be updated, only the content.
- [ ] Develop API to soft delete an evidence entry.
  - [ ] Insert an audit log entry of the delete action.
- [ ] Develop API to hard delete an evidence entry with confirmation steps.
  - [ ] Prompt user for confirmation.
  - [ ] Validate the evidence exists and user permissions.
  - [ ] Log the deletion for auditing purposes.

## Analysis and Extraction

- [ ] Develop API to extract and return the top 10 most used words in all text-based evidence.
  - [ ] Ignore stop words.
- [ ] Develop API to extract and return any links or URLs mentioned in a case by ID.

## Audit and Reporting

- [ ] Develop API to return admin logs for evidence-related actions.
  - [ ] Include details on who added, updated, or deleted evidence and when.
- [ ] Develop API to return a generated report as a PDF with all case details and evidence.

## Bonus Challenges

### Long Polling for Evidence Hard Delete

- [ ] Develop API endpoint for admins to initiate the hard deletion of evidence.
  - [ ] Accept evidence ID and user authentication details.
- [ ] Develop API endpoint for admins to check deletion status using long polling.
  - [ ] Keep the connection open until deletion is complete or timeout occurs.

### Email Notification System

- [ ] Choose and integrate an email service provider.
- [ ] Develop mechanism to trigger email notifications based on specific events.
  - [ ] New crime incidents reported.
  - [ ] Updates to existing cases.

### Case Commenting

- [ ] Develop API endpoints for adding comments to a specific case.
- [ ] Develop API endpoints for retrieving all comments for a case.
- [ ] Develop API endpoints for deleting comments made by assignees.
- [ ] Ensure comments are timestamped and linked to the user.
- [ ] Validate comment length (5-150 characters).
- [ ] Implement rate limiting (no more than 5 comments per minute).

### Deployment

- [ ] Dockerize the project using Docker and Docker Compose.
- [ ] Deploy the application on a cloud platform.
  - [ ] Submit Dockerfile, deployment scripts, and README for running and deploying the project.
