package main

import (
	"Telegram2VCF/types"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"log"
	"os"
	"time"
)

// Should we download the profile pictures or not?
var downloadPhotos bool

func main() {
	client, err := telegram.ClientFromEnvironment(telegram.Options{
		SessionStorage: &session.FileStorage{
			Path: "session",
		},
	})
	if err != nil {
		log.Fatalln("cannot create telegram client:", err)
	}
	// Get the phone number
	phoneNumber := os.Getenv("PHONE")
	if phoneNumber == "" {
		log.Fatalln("please provide your phone number as \"PHONE\" environment variable")
	}
	downloadPhotos = os.Getenv("DOWNLOAD_PHOTOS") == "true"
	// Run the client to get the contacts
	err = client.Run(context.Background(), func(ctx context.Context) error {
		err := client.Auth().IfNecessary(ctx, auth.NewFlow(
			types.SimpleAuth{PhoneNumber: phoneNumber},
			auth.SendCodeOptions{},
		))
		if err != nil {
			return err
		}

		api := client.API()

		contacts, err := api.ContactsGetContacts(ctx, 0)
		if err != nil {
			return err
		}

		contactsModified, ok := contacts.AsModified()
		if !ok {
			return errors.New("not ok")
		}

		return saveContacts(ctx, api, contactsModified.Users)
	})
	if err != nil {
		log.Fatalln("cannot run the client: ", err)
	}
}

func saveContacts(ctx context.Context, client *tg.Client, users []tg.UserClass) error {
	fmt.Println("Fetching contacts...")
	output, err := os.Create("contacts.vcf")
	if err != nil {
		return err
	}
	defer output.Close()
	dl := downloader.NewDownloader()
	var photoBuffer bytes.Buffer // reuse the buffer
	for _, userContact := range users {
		user, ok := userContact.AsNotEmpty()
		if !ok {
			fmt.Println("not ok in", userContact.String())
			continue
		}
		// Download profile photo if possible
		photoBuffer.Reset() // reset the buffer for new file
		if downloadPhotos && user.Photo != nil {
			if photo, ok := user.Photo.AsNotEmpty(); ok && !photo.HasVideo {
				peer := &tg.InputPeerPhotoFileLocation{
					Peer: &tg.InputPeerUser{
						UserID:     user.ID,
						AccessHash: user.AccessHash,
					},
					PhotoID: photo.PhotoID,
				}
				for { // retry if flood
					_, err = dl.Download(client, peer).Stream(ctx, &photoBuffer)
					if flood, _ := tgerr.FloodWait(ctx, err); flood {
						continue // try again!
					}
					if err != nil {
						fmt.Println("cannot download photo:", err)
					}
					break
				}
			}
		}

		err = types.ContactFromUser(user, photoBuffer.Bytes()).AppendAsVCF(output)
		if err != nil {
			return err
		}

		// Sleep a bit to prevent flood
		if downloadPhotos {
			// I think the rate limit for telegram is 30 messages per second
			// So we also provide some space here and go with 25 messages per second
			time.Sleep(time.Second / 25)
		}
	}
	return nil
}
