// Package command holds the verb handlers and the dispatcher
// that maps a parsed token to a handler.
package command

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"terminal-go/parse"
	"terminal-go/world"
)

const defaultSavePath = "zork.save"

func Run(w *world.World, tok parse.Token, out io.Writer) {
	if tok.Verb == "" {
		return
	}
	if dir, ok := world.DirectionFromWord(tok.Verb); ok {
		move(w, dir, out)
		checkWin(w, out)
		return
	}
	switch tok.Verb {
	case "look":
		if tok.Object == "" {
			DescribeRoom(w, w.Player.Room, out)
		} else {
			examine(w, tok, out)
		}
	case "go":
		goCmd(w, tok, out)
	case "take":
		take(w, tok, out)
	case "drop":
		drop(w, tok, out)
	case "inventory":
		inventory(w, out)
	case "examine":
		examine(w, tok, out)
	case "read":
		read(w, tok, out)
	case "save":
		save(w, tok, out)
	case "load":
		load(w, tok, out)
	case "help":
		help(out)
	case "quit", "exit":
		w.Quit = true
		fmt.Fprintln(out, "Goodbye.")
	default:
		fmt.Fprintf(out, "I don't know how to %q.\n", tok.Verb)
	}
	checkWin(w, out)
}

func DescribeRoom(w *world.World, room world.RoomID, out io.Writer) {
	name := w.Rooms.Name[room]
	fmt.Fprintln(out)
	fmt.Fprintln(out, name)
	fmt.Fprintln(out, strings.Repeat("-", len(name)))
	if world.IsDark(w, room) {
		fmt.Fprintln(out, "It is pitch black. You cannot see anything.")
		return
	}
	fmt.Fprintln(out, w.Rooms.Description[room])
	writeItemList(w, room, out)
	writeExitList(w, room, out)
}

func writeItemList(w *world.World, room world.RoomID, out io.Writer) {
	var names []string
	for index := range w.Items.Name {
		id := world.ItemID(index)
		if w.Items.Location[id] != room {
			continue
		}
		names = append(names, w.Items.Name[id])
	}
	if len(names) == 0 {
		return
	}
	sort.Strings(names)
	fmt.Fprintf(out, "You see: %s.\n", strings.Join(names, ", "))
}

func writeExitList(w *world.World, room world.RoomID, out io.Writer) {
	exits := world.ExitsFrom(w, room)
	if len(exits) == 0 {
		fmt.Fprintln(out, "There are no obvious exits.")
		return
	}
	parts := make([]string, 0, len(exits))
	for _, exit := range exits {
		parts = append(parts, world.DirectionName(exit.Dir))
	}
	sort.Strings(parts)
	fmt.Fprintf(out, "Exits: %s.\n", strings.Join(parts, ", "))
}

func goCmd(w *world.World, tok parse.Token, out io.Writer) {
	if len(tok.Args) == 0 {
		fmt.Fprintln(out, "Go where?")
		return
	}
	dir, ok := world.DirectionFromWord(tok.Args[0])
	if !ok {
		fmt.Fprintf(out, "%q is not a direction.\n", tok.Args[0])
		return
	}
	move(w, dir, out)
}

func move(w *world.World, dir world.Direction, out io.Writer) {
	exit, index, ok := world.ExitFrom(w, w.Player.Room, dir)
	if !ok {
		fmt.Fprintln(out, "You can't go that way.")
		return
	}
	if exit.Locked {
		if exit.KeyItem == world.InvalidItem {
			fmt.Fprintln(out, "The way is locked, and you see no way to open it.")
			return
		}
		keyName := w.Items.Name[exit.KeyItem]
		if !world.IsCarrying(w, keyName) {
			fmt.Fprintf(out, "The way is locked. You need a %s.\n", keyName)
			return
		}
		w.Exits[index].Locked = false
		if _, reverseIndex, ok := world.ExitFrom(w, exit.To, world.Opposite(dir)); ok {
			w.Exits[reverseIndex].Locked = false
		}
		fmt.Fprintf(out, "You unlock the way with the %s.\n", keyName)
	}
	w.Player.Room = exit.To
	DescribeRoom(w, exit.To, out)
}

func take(w *world.World, tok parse.Token, out io.Writer) {
	if tok.Object == "" {
		fmt.Fprintln(out, "Take what?")
		return
	}
	if world.IsDark(w, w.Player.Room) {
		fmt.Fprintln(out, "It is too dark to see what you're grabbing for.")
		return
	}
	id := world.FindItemInRoom(w, w.Player.Room, tok.Object)
	if id == world.InvalidItem {
		if world.FindItemInInventory(w, tok.Object) != world.InvalidItem {
			fmt.Fprintln(out, "You're already carrying it.")
			return
		}
		fmt.Fprintln(out, "You don't see that here.")
		return
	}
	if !world.ItemHasTag(w, id, world.ItemTakeable) {
		fmt.Fprintf(out, "You can't take the %s.\n", w.Items.Name[id])
		return
	}
	w.Items.Location[id] = world.InventoryRoom
	fmt.Fprintf(out, "Taken: %s.\n", w.Items.Name[id])
}

func drop(w *world.World, tok parse.Token, out io.Writer) {
	if tok.Object == "" {
		fmt.Fprintln(out, "Drop what?")
		return
	}
	id := world.FindItemInInventory(w, tok.Object)
	if id == world.InvalidItem {
		fmt.Fprintln(out, "You aren't carrying that.")
		return
	}
	w.Items.Location[id] = w.Player.Room
	fmt.Fprintf(out, "Dropped: %s.\n", w.Items.Name[id])
}

func inventory(w *world.World, out io.Writer) {
	var names []string
	for index := range w.Items.Name {
		if w.Items.Location[index] == world.InventoryRoom {
			names = append(names, w.Items.Name[index])
		}
	}
	if len(names) == 0 {
		fmt.Fprintln(out, "You aren't carrying anything.")
		return
	}
	sort.Strings(names)
	fmt.Fprintf(out, "You are carrying: %s.\n", strings.Join(names, ", "))
}

func examine(w *world.World, tok parse.Token, out io.Writer) {
	if tok.Object == "" {
		fmt.Fprintln(out, "Examine what?")
		return
	}
	id := world.FindItemInInventory(w, tok.Object)
	if id == world.InvalidItem {
		if world.IsDark(w, w.Player.Room) {
			fmt.Fprintln(out, "It is too dark to make out any detail.")
			return
		}
		id = world.FindItemInRoom(w, w.Player.Room, tok.Object)
	}
	if id == world.InvalidItem {
		fmt.Fprintln(out, "You don't see that here.")
		return
	}
	fmt.Fprintln(out, w.Items.Description[id])
}

func read(w *world.World, tok parse.Token, out io.Writer) {
	if tok.Object == "" {
		fmt.Fprintln(out, "Read what?")
		return
	}
	id := world.FindItemInInventory(w, tok.Object)
	if id == world.InvalidItem {
		if world.IsDark(w, w.Player.Room) {
			fmt.Fprintln(out, "It is too dark to read.")
			return
		}
		id = world.FindItemInRoom(w, w.Player.Room, tok.Object)
	}
	if id == world.InvalidItem {
		fmt.Fprintln(out, "You don't see that here.")
		return
	}
	if !world.ItemHasTag(w, id, world.ItemReadable) {
		fmt.Fprintf(out, "There is nothing to read on the %s.\n", w.Items.Name[id])
		return
	}
	fmt.Fprintln(out, w.Items.ReadText[id])
}

func save(w *world.World, tok parse.Token, out io.Writer) {
	path := tok.RawObject
	if path == "" {
		path = defaultSavePath
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(out, "Couldn't save: %v\n", err)
		return
	}
	defer file.Close()
	if err := world.Encode(w, file); err != nil {
		fmt.Fprintf(out, "Couldn't save: %v\n", err)
		return
	}
	fmt.Fprintf(out, "Game saved to %s.\n", path)
}

func load(w *world.World, tok parse.Token, out io.Writer) {
	path := tok.RawObject
	if path == "" {
		path = defaultSavePath
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(out, "Couldn't load: %v\n", err)
		return
	}
	defer file.Close()
	loaded, err := world.Decode(file)
	if err != nil {
		fmt.Fprintf(out, "Couldn't load: %v\n", err)
		return
	}
	*w = *loaded
	fmt.Fprintln(out, "Game restored.")
	if w.Won {
		fmt.Fprintln(out, "(This save is from a completed run.)")
		return
	}
	DescribeRoom(w, w.Player.Room, out)
}

func help(out io.Writer) {
	lines := []string{
		"Verbs:",
		"  look, l                            describe the current room",
		"  go <dir>, <dir>                    move (north/south/east/west/up/down, n/s/e/w/u/d)",
		"  take <item>, get, grab, pick up    pick something up",
		"  drop <item>, put down              put something down",
		"  inventory, i                       list what you're carrying",
		"  examine <item>, x                  look closely at something",
		"  read <item>                        read writing on something",
		"  save [file], load [file]           save or restore the game",
		"  help, ?                            show this list",
		"  quit, q                            leave the game",
	}
	for _, line := range lines {
		fmt.Fprintln(out, line)
	}
}

func checkWin(w *world.World, out io.Writer) {
	if w.Won || w.Quit {
		return
	}
	if w.Player.Room != w.GoalRoom {
		return
	}
	if w.Items.Location[w.GoalItem] != world.InventoryRoom {
		return
	}
	w.Won = true
	fmt.Fprintln(out)
	fmt.Fprintln(out, "You step back into the entryway with the coin warm in your hand.")
	fmt.Fprintln(out, "The plaque was right. You have what you came for.")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "*** You win. ***")
}
