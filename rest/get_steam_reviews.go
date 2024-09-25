package rest

import (
	"github.com/arelate/gaugin/rest/compton_data"
	"github.com/arelate/gaugin/rest/gaugin_styles"
	"github.com/arelate/southern_light/steam_integration"
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/compton/consts/weight"
	"github.com/boggydigital/compton/elements/els"
	"github.com/boggydigital/compton/elements/flex_items"
	"github.com/boggydigital/compton/elements/fspan"
	"github.com/boggydigital/compton/elements/iframe_expand"
	"github.com/boggydigital/nod"
	"net/http"
	"strconv"
	"time"
)

const longReviewThreshold = 1024

//type steamReview struct {
//	Author                   string
//	Language                 string
//	Created, Updated         string
//	VotedUp                  bool
//	VotesUp, VotesFunny      int
//	SteamPurchase            bool
//	ReceivedForFree          bool
//	WrittenDuringEarlyAccess bool
//	LongReview               bool
//	Review                   template.HTML
//}

func unixDateFormat(d int64) string {
	return time.Unix(d, 0).Format("Jan 2, 2006")
}

func GetSteamReviews(w http.ResponseWriter, r *http.Request) {

	// GET /steam-reviews?id

	id := r.URL.Query().Get("id")

	sar, err := getSteamReviews(http.DefaultClient, id)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	//reviews := make([]steamReview, 0, len(sar.Reviews))
	//sb := &strings.Builder{}
	//srvm := view_models.NewSteamReviews(sar)

	//for _, rev := range sar.Reviews {
	//	reviews = append(reviews,
	//		steamReview{
	//			Author:                   rev.Author.SteamId,
	//			Language:                 rev.Language,
	//			Created:                  unixDateFormat(rev.TimestampCreated),
	//			Updated:                  unixDateFormat(rev.TimestampUpdated),
	//			VotedUp:                  rev.VotedUp,
	//			VotesUp:                  rev.VotesUp,
	//			VotesFunny:               rev.VotesFunny,
	//			SteamPurchase:            rev.SteamPurchase,
	//			ReceivedForFree:          rev.ReceivedForFree,
	//			WrittenDuringEarlyAccess: rev.WrittenDuringEarlyAccess,
	//			LongReview:               len(rev.Review) > longReviewThreshold,
	//			Review:                   template.HTML(rev.Review),
	//		})
	//}

	section := compton_data.SteamReviewsSection
	ifc := iframe_expand.IframeExpandContent(section, compton_data.SectionTitles[section]).
		AppendStyle(gaugin_styles.SteamReviews)

	pageStack := flex_items.FlexItems(ifc, direction.Column)
	ifc.Append(pageStack)

	for ii, review := range sar.Reviews {
		if srf := steamReviewFragment(ifc, review); srf != nil {
			pageStack.Append(srf)
		}
		if ii < len(sar.Reviews)-1 {
			pageStack.Append(els.Hr())
		}
	}

	if err := ifc.WriteContent(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}

	//if err := tmpl.ExecuteTemplate(sb, "steam-reviews-content", srvm); err != nil {
	//	http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//gaugin_middleware.DefaultHeaders(w)
	//
	//if err := app.RenderSection(id, stencil_app.SteamReviewsSection, sb.String(), w); err != nil {
	//	http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	//	return
	//}
}

func steamReviewFragment(r compton.Registrar, review steam_integration.Review) compton.Element {

	container := flex_items.FlexItems(r, direction.Column).RowGap(size.Normal)

	header := flex_items.FlexItems(r, direction.Column).RowGap(size.XXXSmall)

	votedRow := flex_items.FlexItems(r, direction.Row).ColumnGap(size.Small)
	votedEmoji := "👎"
	votedTitle := "Not Recommended"
	votedColor := color.Red
	if review.VotedUp {
		votedEmoji = "👍"
		votedTitle = "Recommended"
		votedColor = color.Green
	}

	votedRow.Append(
		fspan.Text(r, votedTitle).FontWeight(weight.Bolder).ForegroundColor(votedColor),
		fspan.Text(r, votedEmoji))

	authorRow := appendSteamReviewHeadingRow(r, "Author")
	if review.Author.NumGamesOwned > 0 {
		appendSteamReviewPropertyValue(r, authorRow, "Games owned:", strconv.Itoa(review.Author.NumGamesOwned))
	}
	if review.Author.NumReviews > 0 {
		appendSteamReviewPropertyValue(r, authorRow, "Reviews posted:", strconv.Itoa(review.Author.NumReviews))
	}

	datesRow := appendSteamReviewHeadingRow(r, "Review")
	if review.TimestampCreated > 0 {
		appendSteamReviewPropertyValue(r, datesRow, "Cr:", epochDate(review.TimestampCreated))
	}
	if review.TimestampUpdated > 0 {
		appendSteamReviewPropertyValue(r, datesRow, "Upd:", epochDate(review.TimestampUpdated))
	}

	playtimeRow := appendSteamReviewHeadingRow(r, "Playtime")
	if review.Author.PlaytimeAtReview > 0 {
		appendSteamReviewPropertyValue(r, playtimeRow, "At review:", minutesToHours(review.Author.PlaytimeAtReview))
	}
	if review.Author.PlaytimeLastTwoWeeks > 0 {
		appendSteamReviewPropertyValue(r, playtimeRow, "Last 2w:", minutesToHours(review.Author.PlaytimeLastTwoWeeks))
	}
	if review.Author.PlaytimeForever > 0 {
		appendSteamReviewPropertyValue(r, playtimeRow, "Total:", minutesToHours(review.Author.PlaytimeForever))
	}
	if review.Author.DeckPlaytimeAtReview > 0 {
		appendSteamReviewPropertyValue(r, playtimeRow, "Steam Deck:", minutesToHours(review.Author.DeckPlaytimeAtReview))
	}

	noticeRow := appendSteamReviewHeadingRow(r, "")
	if review.PrimarilySteamDeck {
		appendSteamReviewNotice(r, noticeRow, "Primarily Steam Deck")
	}
	if !review.SteamPurchase {
		appendSteamReviewNotice(r, noticeRow, "Not Steam purchase")
	}
	if review.ReceivedForFree {
		appendSteamReviewNotice(r, noticeRow, "Received for free")
	}
	if review.WrittenDuringEarlyAccess {
		appendSteamReviewNotice(r, noticeRow, "Written during Early Access")
	}

	header.Append(votedRow, authorRow, playtimeRow, datesRow)
	if noticeRow.HasChildren() {
		header.Append(noticeRow)
	}
	container.Append(header)

	var reviewContainer compton.Element
	if len(review.Review) > longReviewThreshold {
		dsTitle := fspan.Text(r, "Expand full review").ForegroundColor(color.Blue)
		dsReview := els.Details().AppendSummary(dsTitle)
		container.Append(dsReview)
		reviewContainer = dsReview
	} else {
		reviewContainer = container
	}
	reviewContainer.Append(fspan.Text(r, review.Review))

	votesRow := appendSteamReviewHeadingRow(r, "Votes")
	if review.VotesUp > 0 {
		appendSteamReviewPropertyValue(r, votesRow, "Helpful:", strconv.Itoa(review.VotesUp))
	}
	if review.VotesFunny > 0 {
		appendSteamReviewPropertyValue(r, votesRow, "Funny:", strconv.Itoa(review.VotesFunny))
	}

	if votesRow.HasChildren() {
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

func appendSteamReviewPropertyValue(r compton.Registrar, c compton.Element, p, v string) {
	c.Append(fspan.Text(r, p).FontSize(size.Small).ForegroundColor(color.Gray))
	c.Append(fspan.Text(r, v).FontSize(size.Small))
}

func appendSteamReviewNotice(r compton.Registrar, c compton.Element, n string) {
	notice := fspan.Text(r, n).
		FontWeight(weight.Bolder).
		FontSize(size.Small).
		ForegroundColor(color.Yellow)
	c.Append(notice)
}

func appendSteamReviewHeadingRow(r compton.Registrar, title string) compton.Element {
	row := flex_items.FlexItems(r, direction.Row).ColumnGap(size.XSmall)
	if title != "" {
		row.Append(fspan.Text(r, title).FontSize(size.Small).FontWeight(weight.Bolder))
	}
	return row

}
