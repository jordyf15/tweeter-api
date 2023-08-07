# Tweeter API
Under Development...

## Endpoint Documentation
### Register User
#### Request
Method: `POST`  
Route: `/register`  
Request Body:  
```
{
    fullname: "Fubuki Shirakami",
    username: "fubuki",
    email: "fubuki@gmail.com",
    password: "Password123!"
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: {
        // user data
    },
    meta: {
        access_token: "access token",
        refresh_token: "refresh token",
        expires_at: 1000000,
    }
}
    
```
### Login User
#### Request
Method: `POST`  
Route: `/login`  
Request Body:
```
{
    login: "fubuki", //username or email
    password: "Password123!"
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: {
        // user data
    },
    meta: {
        access_token: "access token",
        refresh_token: "refresh token",
        expires_at: 1000000,     
    }
}
```
### Get Profile
#### Request
Method: `GET`  
Route: `/users/:user_id`  
Query Params: 
```
{
    filter: "tweets",// valid values are tweets, tweets_n_replies, media, likes  
}
```
Request Header:
```
{
    Authorization: "Bearer accesstoken",
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    user: {
        // user data
    },
    content_filter: "tweets", // values are tweets, tweets_n_replies, media, likes,
    content: {
        data: [
        // data of tweets, tweets and replies, media, or likes
        ],
        meta: {
            per_page: 20,
            next: "next page url",
            prev: "previous page url"
        }
    }
}
```
### Edit Profile
#### Request
Method: `PATCH`  
Route: `/users/:user_id`  
Request Header:
```
{
    Authorization: "Bearer accesstoken",
}
```
Request Body:
```
{
    fullname: "shirakami fubuking",
    username: "newfubuki",
    email: "newfubuki@gmail.com",
    description: "We are friends!",
    profile_image: file,
    background_image: file,
    is_remove_profile_image: true,
    is_remove_background_image: true
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    user: {
        // user data
    }
}
```
### Change Password
#### Request
Method: `POST`  
Route: `/users/:user_id/password/change`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Request Body:
```
{
    old_password: "Password123!",
    new_password: "Password321!"
}
```
#### Response
Status Code: `204`  
### Get User Tweets
#### Request
Method: `GET`  
Route: `/users/:user_id/tweets`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001,
    created_after: 01-01-2002,
    per_page: 20
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [], // user tweets
    meta: {
        next: "next page url",
        prev: "previous page url"
    }
}
```
### Get User Tweets and Replies
#### Request
Method: `GET`  
Route: `/users/:user_id/tweets-and-replies`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001,
    created_after: 01-01-2002,
    per_page: 20
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [], // user tweets and tweets with user replies
    meta: {
        next: "next page url",
        prev: "previous page url"
    }
}
```
### Get User Media
#### Request
Method: `GET`  
Route: `/users/:user_id/media`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001,
    created_after: 01-01-2002,
    per_page: 20
}
```
#### Response
Status Code: `200`  
Response Body: 
```
{
    data: [], // user tweets with media
    meta: {
        next: "next page url",
        prev: "previous page url"
    }
}
```
### Get User Likes
#### Request
Method: `GET`  
Route: `/users/:user_id/likes`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001,
    created_after: 01-01-2002,
    per_page: 20
}
```
#### Response
Status Code: `200`
Response Body:
```
{
    data: [], // tweets that user likes
    meta: {
        next: "next page url",
        prev: "previous page url"
    }
}
```
### Get User Following
#### Request
Method: `GET`  
Route: `/users/:user_id/following`  
Request Header:  
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001,
    created_after: 01-01-2002,
    per_page: 20
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [], // users that is followed by this user
    meta: {
        next: "next page url",
        prev: "previous page url"
    }
}
```
### Get User Follower
#### Request
Method: `GET`  
Route: `/users/:user_id/follower`  
Request Header:  
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001,
    created_after: 01-01-2002,
    per_page: 20
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [], // users that is following this user
    meta: {
        next: "next page url",
        prev: "previous page url"
    }
}
```
### Follow User
#### Request
Method: `POST`  
Route: `/users/{id}/follow`  
Request Header:  
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`  
### Unfollow User
#### Request
Method: `DELETE`  
Route: `/users/:user_id/follow`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`
### Post Tweet
#### Request
Method: `POST`  
Route: `/tweets`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Request Body:
```
{
    description: "description of tweet",
    images: [image1, image2],
    hashtags: ["tag1", "tag2"],
    visibility: "public", // valid values are group or public
    reply_constraint: "everyone", // valid values are everyone or following
    group_id: "id of group tweet is posted", // [optional] only is visibility is group
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    // tweet data
}
```
### Edit Tweet
#### Request
Method: `PATCH`  
Route: `/tweets/:tweet_id`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Request Body:
```
{
    description: "updated description of tweet",
    new_images: [img3, img4],
    deleted_image_ids: ["img1", "img2"],
    reply_constraint: "everyone", // valid values are everyone or following,
    hashtags: ["tag1", "tag2"],
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    // tweet data
}

```
### Get Comments on Tweet
#### Request
Method: `GET`  
Route: `/tweets/:tweet_id/comments`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001,
    created_after: 01-01-2002,
    per_page: 20
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [
        tweet's comments
    ],
    meta: {
        prev: "previous page url",
        next: "next page url"
    }
}
```
### Post Comment on Tweet
#### Request
Method: `POST`  
Route: `/tweets/:tweet_id/comments`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Request Body:
```
{
    comment: "comment",
    images: [img1, img2],
    hashtags: ["tag1", "tag2"]
}
```
#### Response
Status Code: `200`  
Response Body: 
```
{
    // comment data
}
```
### Update Comment on Tweet
#### Request
Method: `PATCH`  
Route: `/tweets/:tweet_id/comments/:comment_id`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Request Body:
```
{
    comment: "comment",
    new_images: [img1, img2],
    hashtags: ["tag1", "tag2"],
    deleted_image_ids: ["img_id1", "img_id2"]
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    // comment data
}
```
### Delete Comment on Tweet
#### Request
Method: `DELETE`  
Route: `/tweets/:tweet_id/comments/:comment_id`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`
### Retweet a Tweet
#### Request
Method: `POST`  
Route: `/tweets/:tweet_id/retweets`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`
### Unretweet a Tweet
#### Request
Method: `DELETE`  
Route: `/tweets/:tweet_id/retweets`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`
### Like a Tweet
#### Request
Method: `POST`  
Route: `/tweets/:tweet_id/likes`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`
### Unlike a Tweet
#### Request
Method: `DELETE`  
Route: `tweets/:tweet_id/likes`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`
### Save a Tweet
#### Request
Method: `POST`
Route: `/tweets/:tweet_id/saves`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`
### Unsave a Tweet
#### Request
Method: `DELETE`  
Route: `/tweets/:tweet_id/saves`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`
### Like a Comment
#### Request
Method: `POST`  
Route: `/tweets/:tweet_id/comments/:comment_id/likes`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`
### Unlike a Comment
Method: `DELETE`  
Route: `/tweets/:tweet_id/comments/:comment_id/likes`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`
### Get Home
#### Request
Method: `GET`  
Route: `/home`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    per_page: 20
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: {
        tweets: [], // tweets of followed users
        trends: [], //hashtags
        follow_recommendation: [], // recommended users to follow
    },
    meta:{
        prev: "previous page url",
        next: "next page url"
    }
}
```
### Get Home Tweets
#### Request
Method: `GET`  
Route: `/home/tweets`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    per_page: 20,
    created_before: 01-01-2001,
    created_after: 01-01-2001
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [
        // tweets of followed users
    ],
    meta: {
        prev: "previous page url",
        next: "next page url"
    }
}
```
### Get Hashtag Tweets
#### Request
Method: `GET`  
Route: `/hashtags/:hashtag_id/tweets`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    per_page: 20,
    created_before: 01-01-2001,
    created_after: 01-01-2001
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [
        // tweets containing the hashtag
    ],
    meta: {
        prev: "previous page url",
        next: "next page url"
    }
}
```
### Get Explore Tweets
#### Request
Method: `GET`  
Route: `/explore/tweets`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001, // for latest and media filter
    created_after: 01-01-2002, // for latest and media filter
    page: 1, // for top filter
    per_page: 20,
    search: "hololive",
    filter: "top" // valid values are top, latest, media
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [
        // tweets
    ],
    meta: {
        prev: "previous page url",
        next: "next page url"
    }
}
```
### Get Explore People
#### Request
Method: `GET`  
Route: `/explore/people`  
Request Header:  
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    page: 1,
    per_page: 20,
    search: "fubuki"
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [
        // popular users
    ],
    meta: {
        prev: "previous page url",
        next: "next page url"
    }
}
```
### Get Bookmark
#### Request
Method: `GET`  
Route: `/bookmark`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001,
    created_after: 02-01-2001,
    per_page: 20,
}
```
#### Response
Status Code: `200`  
Response Body: 
```
{
    data: [
        // bookmarked tweets
    ],
    meta: {
        prev: "previous page url",
        next: "next page url"
    }
}
```
### Get Groups
#### Request
Method: `GET`  
Route: `/groups`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001,
    created_after: 02-01-2001,
    per_page: 20,
    query: "hololive"
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [
        // groups
    ],
    meta: {
        prev: "previous page url",
        next: "next page url"
    }
}
```
### Get Group Tweets
#### Request
Method: `GET`  
Route: `/groups/:group_id/tweets`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001,
    created_after: 02-01-2001,
    per_page: 20
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [
        // tweets
    ],
    meta: {
        prev: "previous page url",
        next: "next page url"
    }
}
```
### Get Group Moderators
#### Request
Method: `GET`  
Route: `/groups/:group_id/moderators`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001,
    created_after: 02-01-2001,
    per_page: 20
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [
        // users that are moderators
    ],
    meta: {
        prev: "previous page url",
        next: "next page url"
    }
}
```
### Get Group Members
#### Request
Method: `GET`  
Route: `/groups/:group_id/members`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Query Params:
```
{
    created_before: 01-01-2001,
    created_after: 02-01-2001,
    per_page: 20
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    data: [
        // users that are members
    ],
    meta: {
        prev: "previous page url",
        next: "next page url"
    }
}
```
### Create Group
#### Request
Method: `POST`  
Route: `/groups`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Request Body:
```
{
    name: "hololive",
    description: "vtuber comedian group hololive",
    image: img1,
}
```
#### Response
Status Code: `200`  
Response Body:
```
{
    // group data
}
```
### Edit Group
#### Request
Method: `PATCH`  
Route: `/groups/:group_id`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
Request Body:
```
{
    name: "hololive",
    description: "vtuber comedian group hololive",
    image: img1,
}
```
#### Response
Status Code: `200`
```
{
    // group data
}
```
### Delete Group
#### Request
Method: `DELETE`  
Route: `/groups/:group_id`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`
### Join Group
#### Request
Method: `POST`  
Route: `/groups/:group_id/members`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`
### Edit Group Member Role
#### Request
Method: `PATCH`  
Route: `/groups/:group_id/members/:member_id`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`  
Response Body:
```
{
    //groupmember data
}
```
### Remove Group Member
#### Request
Method: `DELETE`  
Route: `/groups/:group_id/members/:member_id`  
Request Header:
```
{
    Authorization: "Bearer accesstoken"
}
```
#### Response
Status Code: `204`