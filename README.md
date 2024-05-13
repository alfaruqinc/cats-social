# Cats Social

CatsSocial is a unique application designed to connect cat owners in a fun and interactive way. With CatsSocial, cat enthusiasts can create profiles for their feline friends and discover other cats in their area. The core functionality of the application revolves around matching cats based on various traits, personalities, and interests, fostering new connections within the feline community.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development.

## MakeFile

build the application
```bash
make build
```

run the application
```bash
make run
```

live reload the application
```bash
make watch
```

clean up binary from the last build
```bash
make clean
```

## API

### Authentication

#### Register User
- **Method:** `POST`
- **Endpoint:** `/v1/user/register`
- **Description:** Registers a new user.
- **Request Body:**
  - `name` (string, required): The name of the user.
  - `email` (string, required): The email address of the user.
  - `password` (string, required): The password of the user.
- **Response:** Returns user details upon successful registration.

#### User Login
- **Method:** `POST`
- **Endpoint:** `/v1/user/login`
- **Description:** Logs in an existing user.
- **Request Body:**
  - `email` (string, required): The email address of the user.
  - `password` (string, required): The password of the user.
- **Response:** Returns authentication token upon successful login.

### Manage Cats

#### Create Cat
- **Method:** `POST`
- **Endpoint:** `/v1/cat`
- **Description:** Creates a new cat profile for the authenticated user.
- **Request Body:**
  - `name` (string, required): The name of the cat.
  - `race` (string): The breed of the cat.
  - `sex` (integer): The age of the cat.
  - `ageInMonth` (string): The age of the cat.
  - `description` (string): The description of the cat.
  - `imageUrls` (array of url): The images of the cat.
- **Response:** Returns details of the created cat profile.

#### Get Cats
- **Method:** `GET`
- **Endpoint:** `/v1/cat`
- **Description:** Retrieves all cat profiles
- **Response:** Returns a list of cat profiles.

#### Update Cat
- **Method:** `PUT`
- **Endpoint:** `/v1/cat/{id}`
- **Description:** Updates the details of a cat profile.
- **Request Body:** Same as Create Cat.
- **Response:** Returns updated details of the cat profile.

#### Delete Cat
- **Method:** `DELETE`
- **Endpoint:** `/v1/cat/{id}`
- **Description:** Deletes a cat profile.
- **Response:** Returns a success message upon successful deletion.

### Match Cat

#### Match Cats
- **Method:** `POST`
- **Endpoint:** `/v1/cat/match`
- **Description:** Matches the authenticated user's cat with other cats in the system.
- **Request Body:**
  - `matchCatId` (string, required): The ID of the user's cat to match.
  - `userCatId` (string, required): The ID of the user's cat owner.
  - `message` (string, required): The message.
- **Response:** Returns a success message upon match.

#### Get Matches
- **Method:** `GET`
- **Endpoint:** `/v1/cat/match`
- **Description:** Retrieves all matches for the authenticated user's cat.
- **Response:** Returns a list of matched cats.

#### Approve Match
- **Method:** `POST`
- **Endpoint:** `/v1/cat/match/approve`
- **Description:** Approves a match between the authenticated user's cat and another cat.
- **Request Body:**
  - `matchId` (string, required): The ID of the match to approve.
- **Response:** Returns a success message upon approval.

#### Reject Match
- **Method:** `POST`
- **Endpoint:** `/v1/cat/match/reject`
- **Description:** Rejects a match between the authenticated user's cat and another cat.
- **Request Body:**
  - `matchId` (string, required): The ID of the match to reject.
- **Response:** Returns a success message upon rejection.

#### Delete Match
- **Method:** `DELETE`
- **Endpoint:** `/v1/cat/match/{id}`
- **Description:** Deletes a match between the authenticated user's cat and another cat.
- **Response:** Returns a success message upon deletion.
