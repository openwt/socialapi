package twitter

import (
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/yageek/socialios/social"
)

var api *anaconda.TwitterApi

type Tweet anaconda.Tweet

func init() {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))

	api = anaconda.NewTwitterApi("", "")
}

func Search(word string, lastTweetId string) []social.Data {

	if word == "" {
		return nil
	}
	v := url.Values{}
	v.Set("count", "30")
	if lastTweetId != "" {
		v.Set("max_id", lastTweetId)
	}
	searchResult, _ := api.GetSearch(word, v)

	var tweets []social.Data

	for _, elem := range searchResult.Statuses {
		log.Println("Content:", elem.Text)
		socialData := social.Data{Id: elem.Id, Author: elem.User.Name, Content: elem.Text}

		for _, elem := range elem.Entities.Media {
			socialData.Images = append(socialData.Images, elem.Media_url)
		}
		tweets = append(tweets, socialData)
	}
	return tweets
}

func (t Tweet) Author() string {
	return t.User.Name
}

func (t Tweet) Content() string {
	return t.Text
}

func (t Tweet) Image() url.URL {
	return url.URL{}
}
