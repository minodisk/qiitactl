package command

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/model"
)

func CreateItem(c *cli.Context) {
	err := func() (err error) {
		// client := api.NewClient()
		filename := c.String("file")
		_, err = model.NewItemFromFile(filename)
		if err != nil {
			return
		}
		return
	}()
	if err != nil {
		panic(err)
	}
}

func ShowItems(c *cli.Context) {
	err := func() (err error) {
		client := api.NewClient()
		err = model.ShowItems(client)
		return
	}()
	if err != nil {
		log.Fatal(err)
	}
}

func ShowItem(c *cli.Context) {
}

func FetchItem(c *cli.Context) {
	// Write your code here
}

func FetchItems(c *cli.Context) {
	err := func() (err error) {
		client := api.NewClient()

		items, err := model.FetchItems(client)
		if err != nil {
			return
		}
		err = items.SaveToLocal("mine")
		if err != nil {
			return
		}

		var teams model.Teams
		err = teams.Fetch(client)
		if err != nil {
			return
		}
		for _, team := range teams {
			var items model.Items
			items, err = model.FetchItemsInTeam(client, team)
			if err != nil {
				return
			}
			err = items.SaveToLocal(team.Name)
			if err != nil {
				return
			}
		}
		return
	}()
	if err != nil {
		panic(err)
	}
}

func UpdateItem(c *cli.Context) {
	// Write your code here
}

func UpdateItems(c *cli.Context) {
	// Write your code here
}

func DeleteItem(c *cli.Context) {
	// Write your code here
}

// func fetchAllItems() (err error) {
// 	err = fetchItems("")
// 	if err != nil {
// 		return
// 	}
// 	err = fetchItems("")
// 	return
// }

// func ItemsDiff(commit1, commit2 string) (err error) {
// 	fmt.Printf("Item diff between %s and %s\n", commit1, commit2)
//
// 	err = exec.Command("git", "config", "--local", "core.quotepath", "false").Run()
// 	if err != nil {
// 		return
// 	}
//
// 	cmd := exec.Command("git", "--no-pager", "diff", "--name-only", commit1, commit2)
//
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		return
// 	}
//
// 	filenames, err := cmd.Output()
// 	for _, filename := range strings.Split(string(filenames), "\n") {
// 		if filename == "" {
// 			continue
// 		}
// 		filename = filepath.Join(wd, strings.Trim(filename, "\""))
// 		fmt.Println(filename)
//
// 		b, err := ioutil.ReadFile(filename)
// 		if err != nil {
// 			return err
// 		}
// 		_, err = model.NewItem(string(b))
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	return
// }
