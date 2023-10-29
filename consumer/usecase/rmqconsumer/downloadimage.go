package rmqconsumer

import (
	"consumer/dto"
	"consumer/repository"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strings"

	"github.com/nfnt/resize"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func DownloadImage(m *amqp.Delivery) error {
	data := &dto.RMQProductData{}
	err := json.Unmarshal(m.Body, data)
	if err != nil {
		zap.L().Error("Couldn't Unmarshal data",
			zap.Error(err))
		return err
	}

	// Fetch product from the database using productID
	productImages, err := repository.GetProductFromID(data.Id)
	if err != nil {
		zap.L().Error("Unable to fetch data from DB",
			zap.Error(err))
		return err
	}
	productImagesArray := strings.Split(productImages, ",")

	// Fetch product images, compress, and store locally
	err = saveImage(productImagesArray, data.Id)
	if err != nil {
		zap.L().Error("error saving image to local",
			zap.Error(err))
		return err
	}

	// Update the product's `compressed_product_images` column
	path := generateLocalPaths(data.Id, len(productImagesArray))
	err = repository.UpdateProductImage(data.Id, path)
	if err != nil {
		zap.L().Error("Unable to fetch data from DB",
			zap.Error(err))
		return err
	}

	return nil
}

func saveImage(productImages []string, productID int) error {
	for i, productImage := range productImages {

		outputFilePath := fmt.Sprintf("product_%d_image_%d.jpg", productID, i)
		// Send an HTTP GET request to the URL
		response, err := http.Get(productImage)
		if err != nil {
			fmt.Println("Error while making the request:", err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Printf("HTTP request failed with status: %s\n", response.Status)
			continue
		}

		img, _, err := image.Decode(response.Body)
		if err != nil {
			return err
		}

		img = resize.Resize(800, 0, img, resize.Lanczos3)

		outFile, err := os.Create(outputFilePath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		err = jpeg.Encode(outFile, img, nil)
		if err != nil {
			return err
		}

		fmt.Println("Image downloaded and saved to:", outputFilePath)
	}
	return nil
}

func generateLocalPaths(productID int, numImages int) string {
	paths := make([]string, numImages)
	for i := 0; i < numImages; i++ {
		paths[i] = fmt.Sprintf("product_%d_image_%d.jpg", productID, i)
	}
	return strings.Join(paths, ",")
}
