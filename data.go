package main

import (
	"github.com/charmbracelet/lipgloss"
)

var titles = []string{
	"Study Finds Majority of People Would Rather Be Anywhere Else",
	"New Study Reveals Earth Actually Flat, Scientists Baffled",
	"Nationwide Survey Confirms 97% of Americans Would Trade Everything for a Nap",
	"Local Man Achieves Lifetime Dream of Winning Argument on the Internet",
	"Report: 99% of Parents Consider Going into Witness Protection after Summer Vacation",
	"Study: Humans Prefer Pets Over People, Cats Over Everything",
	"Area Woman Masterfully Avoids Eye Contact for Entire Office Meeting",
	"Breaking: Research Shows Existential Crisis Just a Part of Daily Routine",
	"Survey Shows Millennials Would Rather Die than Answer Phone Call",
	"Scientists Discover Brain Only Capable of Holding Song Lyrics and Useless Trivia",
	"Exclusive: World's Laziest Man Invents a Way to Nap While Sleeping",
	"New Study Finds Procrastination Actually Productive... Eventually",
	"Report: 87% of Adults Still Not Sure What They Want to Be When They Grow Up",
	"Local Woman's Only Form of Exercise Now Justifying Poor Life Choices",
	"Man Brilliantly Solves Global Warming Crisis by Turning Up the AC",
	"Study: 100% of People Would Rather Binge-Watch Netflix Than Socialize",
	"Nationwide Survey Confirms 99% of People's Best Ideas Happen in the Shower",
	"Research Finds Being an Adult Just Endless Cycle of Wanting to Nap",
	"New Study Shows People More Likely to Believe Fake News if It Confirms Existing Beliefs",
	"Local Man Declares Himself Mayor of Couch, Establishes Own Tax System",
}

var surnames = []string{
	"Smith",
	"Johnson",
	"Williams",
	"Jones",
	"Brown",
	"Davis",
	"Miller",
	"Wilson",
	"Moore",
	"Taylor",
	"Anderson",
	"Thomas",
	"Jackson",
	"White",
	"Harris",
	"Martin",
	"Thompson",
	"Garcia",
	"Martinez",
	"Robinson",
}

var names = []string{
	"James",
	"Mary",
	"Robert",
	"Patricia",
	"John",
	"Jennifer",
	"Michael",
	"Linda",
	"David",
	"Elizabeth",
	"William",
	"Barbara",
	"Richard",
	"Susan",
	"Joseph",
	"Jessica",
	"Thomas",
	"Sarah",
	"Christopher",
	"Karen",
}

var emails = []string{
	"skyler.banana@email.com",
	"cosmicjellybean@email.com",
	"pixelprincess@email.com",
	"galaxygobbler@email.com",
	"techno-unicorn@email.com",
	"quantumquasar@email.com",
	"lunarleprechaun@email.com",
	"stardustsurfer@email.com",
	"thundercloud9@email.com",
	"neonnarwhal@email.com",
	"dreamydolphin@email.com",
	"sushisamurai@email.com",
	"velvetvolcano@email.com",
	"auroraborealis@email.com",
	"tangerinetornado@email.com",
	"moonbeam.mermaid@email.com",
	"electricelliot@email.com",
	"enchantedelk@email.com",
	"midnightmystic@email.com",
	"flamingo.fantasy@email.co",
}

var content = []string{
	"Lorem ipsum, dolor sit amet consectetur adipisicing elit. Quaerat, hic eius",
	"distinctio nihil tempore quibusdam temporibus aperiam libero, accusantium",
	"sapiente culpa amet atque, quos molestias delectus corrupti reiciendis. Quam",
	"ipsa voluptate numquam laboriosam exercitationem quos, nisi eligendi vitae",
	"tempora molestias maxime. Quis in fugiat eveniet debitis provident, veritatis",
	"nulla quas ex obcaecati! Id dolores commodi porro odit quae, molestiae dicta",
	"iste nihil veritatis explicabo placeat? Aliquam iure, dolore vitae consequatur",
	"beatae voluptatem voluptatum autem aliquid, est, fuga id saepe! Accusamus enim",
	"fugit nobis dolor vel repudiandae odit, assumenda voluptate quo eligendi amet",
	"repellat modi quas? Ex porro obcaecati distinctio similique error, inventore",
	"mollitia, recusandae assumenda eius voluptatem non amet. Sequi officiis",
	"asperiores beatae. At, quis nesciunt! Architecto, asperiores veniam a laboriosam",
	"officia fuga mollitia tempora tenetur eius provident facilis consequuntur nisi",
	"ipsa omnis molestiae quisquam quaerat optio repellendus laudantium placeat error",
	"exercitationem accusantium animi. Maxime architecto, numquam alias repellat nam",
	"sed unde quod neque enim quis sequi consectetur perferendis ducimus recusandae",
	"dolorem libero fuga sint aliquam mollitia quaerat qui quas. Quisquam possimus",
	"deleniti eum ut voluptate praesentium dolorum autem reprehenderit! Modi itaque",
	"molestias iusto quos tenetur, consequatur esse iure incidunt. Alias veniam",
	"voluptatum voluptate. Cupiditate, excepturi impedit aperiam fuga culpa debitis",
}

var horizontalImg = "https://usnplash.it/300/500"
var verticalImg = "https://unsplash.it/500/300"
var profileImg = "https://unsplash.it/100/100"
var articleImg = "https://unsplash.it/600/400"
var bannerImg = "https://unsplash.it/600/240"

type Info struct {
	style lipgloss.Style
	text  string
}

var helpInfo = []Info{
	{style: helpHeaderStyle, text: "\nSuported types: "},
	{style: helpStyle, text: "string, number, boolean, img, strSet"},
	{style: helpHeaderStyle, text: "\nImg sizes: "},
	{style: helpStyle, text: "default (300x500), vertical (500x300), profile (100x100), article (600x400), banner (600x240)"},
	{style: helpHeaderStyle, text: "\nPredefined string fields: "},
	{style: helpStyle, text: "name, surname/lastName/lastName, email, title, content, author"},
}
