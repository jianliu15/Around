# Geo-index and Image Recognition based Social Network

### Deploy on Google Kubernetes Engine
- Backend: https://around-75015.appspot.com/api/v1
- Frontend: https://around-75015.appspot.com/

### Structure Diagram

#### Backend
<img src="/images/structure_backend.PNG" width="600x">

#### Frontend
<img src="/images/structure_frontend.PNG" width="500x">

### Description

*Backend:*
- Built a scalable web service in Go to handle posts and deployed to Google Cloud (GKE) for better scaling
- Utilized ElasticSearch (GCE) to provide geo-location based search functions such that users can search nearby posts within a distance (e.g. 200km)
- Used Google Dataflow to implement a daily dump of posts to BigQuery table for offline analysis
- Aggregated the data at the post level and user level to improve the keyword based spam detection (BigQuery)
- Used Google Cloud ML API and Tensorflow to train a face detection model and integrate with the Go service.

*Frontend:*
- Built a geo-based social network web application with React JS
- Implemented basic token based registration/login/logout flow with React Router v4 and server-side user authentication with JWT
- Implemented features such as "Create Post", "Nearby Posts As Gallery" and "Nearby Posts In Map" with Ant Design, GeoLocation API and Google Map API
