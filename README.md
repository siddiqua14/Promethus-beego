
# Booking-Rental-App Project with Beego Framework

The Booking App fetches and stores hotel and rental property data from the Booking.com API into a local PostgreSQL database. It implements a backend API using the Beego framework to serve property data and provides a frontend for displaying property listings and details.
## Backend Project (rental)
### Goals

1. **Database**: 
   - Fetch data from the Booking.com API.
   - Store data in a local PostgreSQL database with the following tables:
     - **Location Table**: Contains location names and their associated hotels.
     - **Rental Property Table**: Stores basic information about properties (name, type, bedrooms, bathrooms, amenities).
     - **Property Details Table**: Stores images, descriptions, and other relevant information about properties.

2. **API**:
    
- **Fetch Locations**: 
  - **Route**: `/fetch_locations`
  - **Controller**: `LocationController`
  - **Method**: `GET`
  - **Description**: Fetch and store locations and associated hotels.

- **Fetch Stays Data**:
  - **Route**: `/fetch_stays_data`
  - **Description**: Fetch filtered stays data based on specified criteria.

- **Fetch Hotel Details**:
  - **Route**: `/fetch-hotel-details`
  - **Description**: Fetch detailed information for a specific hotel.

- **Fetch Hotel Images and Descriptions**:
  - **Route**: `/fetch-hotel-images-and-description`
  - **Description**: Fetch images and descriptions for a hotel.

- **Property Listing**:
  - **Route**: `/v1/property/list`
  - **Description**: Provides a list of properties for a specific location. The data is fetched from the database after calling the backend API to store property data.

- **Property Details**:
  - **Route**: `/v1/property/details`
  - **Description**: Provides detailed information for a specific property.
#### Note: 
The system can use a pre-populated database `BookingCom.sql` for property-related endpoints (`/v1/property/list` and `/v1/property/details`). 
This database contains exported data that can be accessed directly without calling external APIs for storing data into database

## Prerequisites
Before you begin, ensure you have met the following requirements:

- **Go (Golang)**: Version 1.18 or higher.
- **Beego Framework**: Installed globally.
- **Git**: For cloning the repository.

## Installation and Setup

### Step 1: Install Go
   - Download and install Go from [https://go.dev/dl/](https://go.dev/dl/).
   - Verify installation:
     ```bash
     go version
     ```
### Step 2: Clone the repository
To get started, clone the project repository:
1. Navigate to your Go src `/go/src` directory to ensure the project is placed in the correct directory for your Go workspace.
2. Clone the repository:
```bash
git clone https://github.com/siddiqua14/rental.git
cd rental
```


### Step 3: Install Beego Framework
Beego is the framework used for this project, and Bee CLI is a development tool. If `Beego` is not installed in your workspace,
- Install them by running:
```bash
go get github.com/beego/beego/v2@latest
```
Ensure your `GOPATH` is set up correctly:
##### To verify installation:
```bash
bee version
```

### Step 4: Docker
Set up the PostgreSQL database using Docker Compose:
- Ensure Docker and Docker Compose are installed on your machine.
- Run the following command to start the database:
```bash
docker-compose up 
```
- This command will create and start the PostgreSQL database container defined in docker-compose.yml.
- Connect to the PostgreSQL database:
```bash
docker exec -it postgres_booking psql -U postgres
\c booking_db
```
### Step 5: Install Dependencies

After installing the required Go modules and Beego dependencies:

```bash
go mod tidy
```
This will automatically resolve any missing dependencies and update your go.sum file with the required entries.
Run the application:
```bash
bee run 
```
Access the application at `http://localhost:8080`.


## Frontend Project (rental_app)
### Goal:
The frontend is designed to show property listings for a location and property details. The goal is to display properties dynamically using AJAX, with a focus on the Property Listing and Property Details pages.
#### To set up the Frontend Project: 
- Go to this repository [rental_app](https://github.com/siddiqua14/rental_app.git) and follow it's instuctions. 
