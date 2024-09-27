package compton_fragments

import (
	"github.com/arelate/southern_light/steam_integration"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/align"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/font_weight"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/elements/details_summary"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"strconv"
	"time"
)

const longReviewThreshold = 750

func SteamReview(r compton.Registrar, review steam_integration.Review) compton.Element {

	container := flex_items.FlexItems(r, direction.Column).RowGap(size.Normal)

	//container.Append(els.Text(strconv.Itoa(len(review.Review))))

	votedTitle := "Not Recommended"
	votedColor := color.Red
	if review.VotedUp {
		votedTitle = "Recommended"
		votedColor = color.Green
	}

	votedHeading := els.H3()
	votedHeading.Append(fspan.Text(r, votedTitle).ForegroundColor(votedColor))

	container.Append(votedHeading)

	header := flex_items.FlexItems(r, direction.Row).ColumnGap(size.Small).RowGap(size.Unset)

	authorRow := SteamReviewHeadingRow(r, "Author")
	if review.Author.NumGamesOwned > 0 {
		AppendSteamReviewPropertyValue(r, authorRow, "Games:", strconv.Itoa(review.Author.NumGamesOwned))
	}
	if review.Author.NumReviews > 0 {
		AppendSteamReviewPropertyValue(r, authorRow, "Reviews:", strconv.Itoa(review.Author.NumReviews))
	}

	datesRow := SteamReviewHeadingRow(r, "Review")
	if review.TimestampCreated > 0 {
		AppendSteamReviewPropertyValue(r, datesRow, "Cr:", epochDate(review.TimestampCreated))
	}
	if review.TimestampUpdated > 0 {
		AppendSteamReviewPropertyValue(r, datesRow, "Upd:", epochDate(review.TimestampUpdated))
	}

	playtimeRow := SteamReviewHeadingRow(r, "Playtime")
	if review.Author.PlaytimeAtReview > 0 {
		AppendSteamReviewPropertyValue(r, playtimeRow, "At review:", minutesToHours(review.Author.PlaytimeAtReview))
	}
	if review.Author.PlaytimeLastTwoWeeks > 0 {
		AppendSteamReviewPropertyValue(r, playtimeRow, "Last 2w:", minutesToHours(review.Author.PlaytimeLastTwoWeeks))
	}
	if review.Author.PlaytimeForever > 0 {
		AppendSteamReviewPropertyValue(r, playtimeRow, "Total:", minutesToHours(review.Author.PlaytimeForever))
	}
	if review.Author.DeckPlaytimeAtReview > 0 {
		AppendSteamReviewPropertyValue(r, playtimeRow, "Steam Deck:", minutesToHours(review.Author.DeckPlaytimeAtReview))
	}

	noticeRow := SteamReviewHeadingRow(r, "")
	if review.PrimarilySteamDeck {
		AppendSteamReviewNotice(r, noticeRow, "Primarily Steam Deck")
	}
	if !review.SteamPurchase {
		AppendSteamReviewNotice(r, noticeRow, "Not Steam purchase")
	}
	if review.ReceivedForFree {
		AppendSteamReviewNotice(r, noticeRow, "Received for free")
	}
	if review.WrittenDuringEarlyAccess {
		AppendSteamReviewNotice(r, noticeRow, "Written during Early Access")
	}

	header.Append(authorRow, playtimeRow, datesRow)
	if noticeRow.HasChildren() {
		header.Append(noticeRow)
	}
	container.Append(header)

	var reviewContainer compton.Element
	if len(review.Review) > longReviewThreshold {
		dsTitle := fspan.Text(r, "Show full review").
			ForegroundColor(color.Blue).
			FontWeight(font_weight.Bolder)
		dsReview := details_summary.Smaller(r, dsTitle, false)
		container.Append(dsReview)
		reviewContainer = dsReview
	} else {
		reviewContainer = container
	}
	reviewContainer.Append(fspan.Text(r, review.Review))

	votesRow := SteamReviewHeadingRow(r, "Votes")
	if review.VotesUp > 0 {
		AppendSteamReviewPropertyValue(r, votesRow, "Helpful:", strconv.Itoa(review.VotesUp))
	}
	if review.VotesFunny > 0 {
		AppendSteamReviewPropertyValue(r, votesRow, "Funny:", strconv.Itoa(review.VotesFunny))
	}

	if review.VotesUp > 0 || review.VotesFunny > 0 {
		container.Append(votesRow)
	}

	return container
}

func minutesToHours(m int) string {
	return strconv.Itoa(m/60) + "h"
}

func epochDate(e int64) string {
	return time.Unix(e, 0).Format("Jan 2, '06")
}

func AppendSteamReviewPropertyValue(r compton.Registrar, c compton.Element, p, v string) {
	c.Append(fspan.Text(r, p).FontSize(size.Small).ForegroundColor(color.Gray))
	c.Append(fspan.Text(r, v).FontSize(size.Small))
}

func AppendSteamReviewNotice(r compton.Registrar, c compton.Element, n string) {
	notice := fspan.Text(r, n).
		FontWeight(font_weight.Bolder).
		FontSize(size.Small).
		ForegroundColor(color.Yellow)
	c.Append(notice)
}

func SteamReviewHeadingRow(r compton.Registrar, title string) compton.Element {
	row := flex_items.FlexItems(r, direction.Row).ColumnGap(size.XSmall).RowGap(size.Unset).AlignItems(align.Center)
	if title != "" {
		row.Append(fspan.Text(r, title).FontSize(size.Small).FontWeight(font_weight.Bolder))
	}
	return row

}
