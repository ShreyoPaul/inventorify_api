package routes

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func Auth(r *gin.Engine) {
	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.Signup)
}

func Inventory(r *gin.Engine) {
	r.GET("/", controllers.GetAllInv)
	r.OPTIONS("/", controllers.GetAllInv)
	r.POST("/", controllers.CreateInv)
	r.POST("/add-attributes-2-inventory", controllers.AddAttributes)
	r.DELETE("/inv/:inventory", controllers.DeleteInvenroty)
}

func Product(r *gin.Engine) {
	r.GET("/:inventory", controllers.GetAllProducts)
	r.GET("/pipeline", controllers.PipelineAllProducts)
	r.POST("/:inventory", controllers.CreateProduct)
	r.PATCH("/", controllers.UpdateProduct)
	r.DELETE("/:inventory", controllers.DeleteProduct)
	r.PATCH("/bill", controllers.BillingProduct)
}

func Billing(r *gin.Engine) {
	r.POST("/billing", controllers.StoreBill)
	r.GET("/billing", controllers.GetAllBills)
}
