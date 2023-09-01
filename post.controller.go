package controllers

import (
	"net/http"
	"strconv"
  	"strings"
//      "net/smtp"
//	"os"
//	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-mongodb/models"
	"github.com/wpcodevo/golang-mongodb/services"
)

type PostController struct {
	postService services.PostService
userService services.UserService
}

func NewPostController(postService services.PostService, userService services.UserService) PostController {
return PostController{postService, userService}
}
func (pc *PostController) CreatePost(ctx *gin.Context) {
	var post *models.CreatePostRequest

	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newPost, err := pc.postService.CreatePost(post)

	if err != nil {
		// if strings.Contains(err.Error(), "title already exists") {
		// 	ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
		// 	return
		// }

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newPost})
}
func (pc *PostController) Unlikepost(ctx *gin.Context) {
	var follow *models.CreatelikeRequest
	if err := ctx.ShouldBindJSON(&follow); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := pc.postService.DeleteLike(follow.Userid, follow.Postid, follow.Like)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Unliked Successfully"})
}
func (pc *PostController) Commentpost(ctx *gin.Context) {
	var comment *models.CreatecommentRequest
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	createcomment, err := pc.postService.CreateComment(comment)
	if err != nil {
		// if strings.Contains(err.Error(), "title already exists") {
		// 	ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
		// 	return
		// }

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": createcomment})

}
func (pc *PostController) UpdatePost(ctx *gin.Context) {
	postId := ctx.Param("postId")

	var post *models.UpdatePost
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedPost, err := pc.postService.UpdatePost(postId, post)
	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedPost})
}
func (pc *PostController) DeletePost(ctx *gin.Context) {
	postId := ctx.Param("postId")

	err := pc.postService.DeletePost(postId)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (pc *PostController) FindPostById(ctx *gin.Context) {
	postId := ctx.Param("postId")

	post, err := pc.postService.FindPostById(postId)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": post})
}
func (pc *PostController) FindPosts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "100")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	posts, err := pc.postService.FindPosts(intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Retrieve user data for each post
                for i := range posts {
		id, err := pc.userService.FindUserById(posts[i].User)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
posts[i].Name = id.Name
		posts[i].User = id.ID.Hex()
		posts[i].Email = id.Email
		posts[i].Location = id.Location
		posts[i].User_group = id.User_group
		posts[i].Profileimg = id.Profileimg
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(posts), "data": posts})
}

func (pc *PostController) FindPostsByUserID(ctx *gin.Context) {
	userID := ctx.Query("userID")

	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "1000")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	posts, err := pc.postService.FindPostsByUserID(userID, intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	for i := range posts {
		id, err := pc.userService.FindUserById(posts[i].User)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		posts[i].Name = id.Name
		posts[i].User = id.ID.Hex()
		posts[i].Email = id.Email
		posts[i].Location = id.Location
		posts[i].User_group = id.User_group
		posts[i].Profileimg = id.Profileimg

	}
	// Retrieve additional user data for each post...

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(posts), "data": posts})
}
func (pc *PostController) Likepost(ctx *gin.Context) {
	var like *models.CreatelikeRequest
	if err := ctx.ShouldBindJSON(&like); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
      	exists, err := pc.postService.ChecklikeExistence(like)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if exists {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Record already exists"})
		return
	}
	createlike, err := pc.postService.CreateLike(like)
	if err != nil {
		// if strings.Contains(err.Error(), "title already exists") {
		// 	ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
		// 	return
		// }

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": createlike})

}
func (pc *PostController) GetLikeByUserID(ctx *gin.Context) {
	userID := ctx.Query("userID")
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "1000")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	posts, err := pc.postService.LikesByUserID(userID, intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	// gtposts, err := pc.postService.FindPostsByUserID(userID, intPage, intLimit)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	// 	return
	// }

	// for j := range gtposts {
	// gpost, err := pc.postService.FindPostsByUserID(posts[j].Userid, intPage, intLimit)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	// 	return
	// }
	// posts[j].Content = gpost.
	// }

	for i := range posts {
		id, err := pc.userService.FindUserById(posts[i].Userid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		gposts, err := pc.postService.FindPostsBypostID(posts[i].Postid, intPage, intLimit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		posts[i].Name = id.Name
		posts[i].Email = id.Email
		posts[i].Location = id.Location
		posts[i].Profileimg = id.Profileimg
		posts[i].Coverimg = id.Coverimg
		posts[i].User_group = id.User_group
		if len(gposts) > 0 {
			posts[i].Content = gposts[0].Content
		}
	}

	// Retrieve additional user data for each post...

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(posts), "data": posts})
}
func (pc *PostController) Followreq(ctx *gin.Context) {
	var follow *models.SendfollowRequest
	if err := ctx.ShouldBindJSON(&follow); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Check if the record already exists
	exists, err := pc.postService.CheckFollowExistence(follow)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if exists {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Record already exists"})
		return
	}

	followReq, err := pc.postService.Followuser(follow)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": followReq})
}
func (pc *PostController) Deletefollowreq(ctx *gin.Context) {
	var follow *models.SendfollowRequest
	if err := ctx.ShouldBindJSON(&follow); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := pc.postService.DeleteFollow(follow.Followbyuser, follow.Followuserid, follow.Status)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data deleted successfully"})
}
func (pc *PostController) FollowreqbyUserID(ctx *gin.Context) {
	userID := ctx.Query("userID")

	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "1000")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	posts, err := pc.postService.FollowByUserID(userID, intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	for i := range posts {
		id, err := pc.userService.FindUserById(posts[i].Followuserid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		posts[i].Name = id.Name
		posts[i].Email = id.Email
		posts[i].Location = id.Location
		posts[i].User_group = id.User_group
		posts[i].Profileimg = id.Profileimg

	}
	// Retrieve additional user data for each post...

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(posts), "data": posts})
}
func (pc *PostController) HandleContactform(ctx *gin.Context) {
	var post *models.ContactModel

	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newPost, err := pc.postService.Contact(post)

	if err != nil {
		log.Println("Failed to save contact record:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save contact record"})
		return
	}
//	if err := sendEmailToAdmin(post); err != nil {
//		log.Println("Failed to send email to admin:", err)
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
//		return
//	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Contact form submitted successfully", "data": newPost})
}
//func sendEmailToAdmin(contact *models.ContactModel) error {
//	from := os.Getenv("SMTP_FROM")
//	password := os.Getenv("SMTP_PASSWORD")
//	to := os.Getenv("SMTP_TO")
//	smtpServer := os.Getenv("SMTP_SERVER")
//	smtpPort := os.Getenv("SMTP_PORT")

//	auth := smtp.PlainAuth("", from, password, smtpServer)

//	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: New Contact Form Submission\n\nName: %s\nEmail: %s\nMessage: %s",
//		from, to, contact.FirstName, contact.Email, contact.Message)

//	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{to}, []byte(message))
//	if err != nil {
//		return err
//	}

//	return nil
//}
func (pc *PostController) Interested(ctx *gin.Context) {
	var interest *models.InterestedRequest
	if err := ctx.ShouldBindJSON(&interest); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Check if the record already exists
	exists, err := pc.postService.CheckinterestExistence(interest)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if exists {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Record already exists"})
		return
	}

	interestReq, err := pc.postService.Interestuser(interest)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": interestReq})

}
func (pc *PostController) Deleteinterestreq(ctx *gin.Context) {
	var interest *models.InterestedRequest
	if err := ctx.ShouldBindJSON(&interest); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := pc.postService.DeleteInterest(interest.Interestbyuser, interest.Interestuserid, interest.Status)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data deleted successfully"})
}
func (pc *PostController) Checkfollowuser(ctx *gin.Context) {
	var follow *models.SendfollowRequest
	if err := ctx.ShouldBindJSON(&follow); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Check if the record already exists
	exists, err := pc.postService.CheckFollowExistence(follow)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if exists {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Record exists"})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Record does not exist"})
	}
}
func (pc *PostController) Checkinterestuser(ctx *gin.Context) {
	var follow *models.InterestedRequest
	if err := ctx.ShouldBindJSON(&follow); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Check if the record already exists
	exists, err := pc.postService.CheckinterestExistence(follow)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if exists {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Record exists"})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Record does not exist"})
	}
}
