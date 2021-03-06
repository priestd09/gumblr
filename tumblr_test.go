package tumblr

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

var testClient *Tumblr

func setup() {
	var credentials map[string]string
	credentialsFile, err := ioutil.ReadFile("./credentials_test.json")
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(credentialsFile, &credentials)
	testClient = New(
		credentials["consumer_key"],
		credentials["consumer_secret"],
		credentials["oauth_key"],
		credentials["oauth_secret"],
	)
}

func TestNew(t *testing.T) {
	var credentials map[string]string
	credentialsFile, err := ioutil.ReadFile("./credentials_test.json")
	if err != nil {
		t.Errorf("error opening the credentials file: %s\n", err)
	}
	json.Unmarshal(credentialsFile, &credentials)
	client := New(
		credentials["consumer_key"],
		credentials["consumer_secret"],
		credentials["oauth_key"],
		credentials["oauth_secret"],
	)
	client.UserInfo() // Without authorization, this command can't be performed
}

func TestBlogInfo(t *testing.T) {
	setup()
	blogInfo := testClient.BlogInfo("staff.tumblr.com")
	if blogInfo.Blog.Name != "staff" {
		t.Error("Client connected to incorrect blog")
	}
}

func TestBlogAvatar(t *testing.T) {
	setup()
	blogAvatar := testClient.BlogAvatar("staff.tumblr.com")
	if reflect.DeepEqual(blogAvatar, make([]byte, 1)) {
		t.Error("Avatar return type is not equivalent")
	}
}

func TestBlogAvatarAndSize(t *testing.T) {
	setup()
	blogAvatar := testClient.BlogAvatarAndSize("staff.tumblr.com", 16)
	if reflect.DeepEqual(blogAvatar, make([]byte, 1)) {
		t.Error("Avatar return type is not equivalent")
	}
}

func TestBlogLikes(t *testing.T) {
	setup()
	params := map[string]string{
		"limit": "20",
	}
	blogLikes := testClient.BlogLikes("mattcunningham.net", params)
	if blogLikes.LikedCount <= 0 {
		t.Error("Incorrect like count returned")
	}
	for _, post := range blogLikes.LikedPost {
		if post.BlogName == "" {
			t.Error("Invalid blog name returned")
		}
	}
}

func TestBlogFollowers(t *testing.T) {
	setup()
	params := map[string]string{
		"limit": "20",
	}
	blogFollowers := testClient.BlogFollowers("mattcunningham.net", params)
	if blogFollowers.TotalUsers <= 0 {
		t.Error("Incorrect follower count returned")
	}
	for _, user := range blogFollowers.Users {
		if user.Name == "" {
			t.Error("Invalid follower name returned")
		}
	}
}

func TestBlogPosts(t *testing.T) {
	setup()
	params := map[string]string{
		"limit": "20",
	}
	blogPosts := testClient.BlogPosts("staff.tumblr.com", params)
	if blogPosts.Blog.Name != "staff" {
		t.Error("Incorrect short blog name")
	}
}

func TestBlogQueuedPosts(t *testing.T) {
	setup()
	params := map[string]string{
		"limit": "20",
	}
	queuedPosts := testClient.BlogQueuedPosts("mattcunningham.net", params)
	for _, post := range queuedPosts.Posts {
		if post.BlogName == "" {
			t.Error("Incorrect short blog name")
		}
	}
}

func TestPost(t *testing.T) {
	setup()
	params := map[string]string{
		"state": "private",
		"type":  "text",
		"title": "Testing Title",
		"body":  "Test text",
	}
	response := testClient.Post("testnames.tumblr.com", params)
	if response.Status != 201 {
		t.Errorf("Test post did not post, response returned %d\n", response.Status)
	}
}

func TestPostEdit(t *testing.T) {
	setup()
	params := map[string]string{
		"state": "private",
		"type":  "text",
		"title": "Testing Title",
		"body":  "Testing text",
	}
	response := testClient.PostEdit("testnames.tumblr.com", 123923316127, params)
	if response.Status != 200 {
		t.Errorf("Test post was not edited, response returned %d\n", response.Status)
	}
}

func TestPostReblog(t *testing.T) {
	setup()
	params := map[string]string{
		"comment": "Test comment",
	}
	response := testClient.PostReblog("testnames.tumblr.com", 122517491420, "kaGXZHdj", params)
	if response.Status != 201 {
		t.Errorf("Test reblog was not reblogged, response returned %d", response.Status)
	}
}

func TestPostDelete(t *testing.T) {
	setup()
	params := map[string]string{
		"limit": "20",
	}
	blogPosts := testClient.BlogPosts("testnames.tumblr.com", params)
	blogPost := blogPosts.Posts[0]
	response := testClient.PostDelete("testnames.tumblr.com", blogPost.ID)
	if response.Status != 200 {
		t.Errorf("Test reblog was not reblogged, response returned %d", response.Status)
	}
}

func TestUserInfo(t *testing.T) {
	setup()
	userInfo := testClient.UserInfo()
	if userInfo.User.Likes <= 0 {
		t.Errorf("User info didn't return the accurate like count")
	}
}

func TestUserDashboard(t *testing.T) {
	setup()
	params := map[string]string{
		"limit": "20",
	}
	blogList := testClient.UserDashboard(params)
	if len(blogList.Posts) <= 0 {
		t.Errorf("User dashboard didn't return the accurate blog post count")
	}
}

func TestUserLikes(t *testing.T) {
	setup()
	params := map[string]string{
		"limit": "20",
	}
	userLikes := testClient.UserLikes(params)
	if userLikes.LikedCount <= 0 {
		t.Errorf("User like count didn't return the accurate count")
	}
}

func TestUserFollowing(t *testing.T) {
	setup()
	params := map[string]string{
		"limit": "20",
	}
	userInfo := testClient.UserFollowing(params)
	if userInfo.TotalBlogs <= 0 {
		t.Errorf("User following didn't return the accurate blog following count")
	}
}

func TestUserFollow(t *testing.T) {
	setup()
	response := testClient.UserFollow("testnames.tumblr.com")
	if response.Status != 200 {
		t.Errorf("Test user was not followed, response returned %d", response.Status)
	}
}

func TestUserUnfollow(t *testing.T) {
	setup()
	response := testClient.UserUnfollow("testnames.tumblr.com")
	if response.Status != 200 {
		t.Errorf("Test user was not unfollowed, response returned %d", response.Status)
	}
}

func TestUserLike(t *testing.T) {
	setup()
	response := testClient.UserLike(122517491420, "kaGXZHdj")
	if response.Status != 200 {
		t.Errorf("Test blog was not liked, response returned %d", response.Status)
	}
}

func TestUserUnlike(t *testing.T) {
	setup()
	response := testClient.UserUnlike(122517491420, "kaGXZHdj")
	if response.Status != 200 {
		t.Errorf("Test blog was not liked, response returned %d", response.Status)
	}
}

func TestTaggedPosts(t *testing.T) {
	setup()
	params := map[string]string{
		"limit": "20",
	}
	taggedPosts := testClient.TaggedPosts("gif", params)
	if len(taggedPosts) <= 0 {
		t.Error("Tagged posts 'gif' did not properly return posts")
	}
}
