package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "os"
    "strings"

    "gopkg.in/yaml.v2"
    "github.com/gempir/go-twitch-irc/v2"
    "github.com/schollz/closestmatch"
)

var keys []string
var keyMatching *closestmatch.ClosestMatch
var weaponsInfo *weapons
var armorInfo *armors
var monstersInfo *monsters
var toolsInfo *tools
var wandsInfo *wands
var ringsInfo *rings
var propsInfo *properties
var comestiblesInfo * comestibles
var potionsInfo *potions
var artifactsInfo *artifacts
var appearsAs *appearances

var allowedBroadcasters *allowedChannels

type allowedChannels struct {
    Names []string `yaml:"channels"`
}

type appearances struct {
    Items map[string]string `yaml:"appearances"`
}

type weapons struct {
    Items map[string]weapon `yaml:"weapons"`
}

type weapon struct {
    Skill string `yaml:"skill"`
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    DamageSmall string `yaml:"damage-small"`
    DamageLarge string `yaml:"damage-large"`
    Material string `yaml:"material"`
}

type armors struct {
    Items map[string]armor `yaml:"armors"`
}

type armor struct {
    Type string `yaml:"skill"`
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    AC int `yaml:"ac"`
    Material string `yaml:"material"`
    Effect string `yaml:"effect"`
    MC int `yaml:"mc"`
}

type artifacts struct {
    Items map[string]artifact `yaml:"artifacts"`
}

type artifact struct {
    BaseItem string `yaml:"base_item"`
    Alignment string `yaml:"alignment"`
    Intelligent bool `yaml:"intelligent"`
    Use string `yaml:"use"`
    Carried string `yaml:"carried"`
    Used string `yaml:"used"`
    Invoked string `yaml:"invoked"`
    Obtaining string `yaml:"obtaining"`
}

type tools struct {
    Items map[string]tool `yaml:"tools"`
}

type tool struct {
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    Use string `yaml:"use"`
    Magic bool `yaml:"magic"`
}

type wands struct {
    Items map[string]wand `yaml:"wands"`
}

type wand struct {
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    Type string `yaml:"type"`
    StartingCharges string `yaml:"starting_charges"`
    Effect string `yaml:"effect"`
    Broken string `yaml:"broken"`
}

type rings struct {
    Items map[string]ring `yaml:"rings"`
}

type ring struct {
    Cost int `yaml:"cost"`
    ExtrinsicGranted string `yaml:"extrinsic_granted"`
    Notes string `yaml:"notes"`
}

type potions struct {
    Items map[string]potion `yaml:"potions"`
}

type potion struct {
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    Effect string `yaml:"effect"`
}

type properties struct {
    Items map[string]property `yaml:"properties"`
}

type property struct {
    Effect string `yaml:"effect"`
    Sources []string `yaml:"sources"`
}

type comestibles struct {
    Items map[string]comestible `yaml:"comestibles"`
}

type comestible struct {
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    NutritionalValue int `yaml:"nutritional_value"`
    Time int `yaml:"time"`
    Conduct string `yaml:"conduct"`
    Effect string `yaml:"effect"`
}


type monsters struct {
    Items map[string]monster `yaml:"monsters"`
}

type monster struct {
    Difficulty int `yaml:"difficulty"`
    Attacks string `yaml:"attacks"`
    Speed int `yaml:"speed"`
    AC int `yaml:"ac"`
    MR int `yaml:"mr"`
    Weight int `yaml:"weight"`
    Alignment string `yaml:"alignment"`
    Genocidable bool `yaml:"genocidable"`
    NutritionalValue int `yaml:"nutritional_value"`
    Size string `yaml:"size"`
    Resistances string `yaml:"resistances"`
    ResistancesConveyed string `yaml:"resistances_conveyed"`
    CorpseSafe bool `yaml:"corpse_safe"`
    Elbereth bool `yaml:"elbereth"`
    Extra string `yaml:"extra"`
}

func getArmorMessage(name string) string {
    var output string 
    if val, ok := armorInfo.Items[name]; ok {
        output = fmt.Sprintf("A %s has AC %d, MC %d, weight %d, costs %dzm and is made of %s. %s",
        strings.ReplaceAll(name,"-"," "), 
        val.AC, 
        val.MC, 
        val.Weight,
        val.Cost,
        val.Material,
        val.Effect)
    }
    return output
}

func getWeaponMessage(name string) string {
    var output string 
    if val, ok := weaponsInfo.Items[name]; ok {
        output = fmt.Sprintf("A %s does %s/%s damage. "+ 
        "It is made of %s, weighs "+ 
        "%d, and is valued at "+
        "%dzm. Works your skill with %s.", 
        name, // strings.ReplaceAll(name,"-"," "), 
        val.DamageSmall, 
        val.DamageLarge, 
        val.Material, 
        val.Weight, 
        val.Cost, 
        val.Skill)
    }
    return output
}

func getArtifactMessage(name string) string {
    var output string 
    if val, ok := artifactsInfo.Items[name]; ok {
        var intelligent string
        if val.Intelligent {
            intelligent = " and intelligent"
        }
        var carried string
        if val.Carried != "" {
            carried = fmt.Sprintf(" While carried it bestows %s.", val.Carried)
        }
        var used string
        if val.Used != "" {
            used = fmt.Sprintf(" While %s it bestows %s.", val.Use, val.Used)
        }
        var invoked string
        if val.Invoked != "" {
            invoked = fmt.Sprintf(" When invoked it bestows %s.", val.Invoked)
        }
        var obtaining string
        if val.Obtaining != "" {
            obtaining = fmt.Sprintf(" %s", val.Obtaining)
        }
        output = fmt.Sprintf("%s is a %s%s artifact whose base item is a %s.%s%s%s%s",
        strings.Title(strings.ReplaceAll(name,"-"," ")), 
        val.Alignment,
        intelligent,
        strings.ReplaceAll(val.BaseItem,"-"," "),
        carried,
        used,
        invoked,
        obtaining)
    }
    return output
}

func getToolMessage(name string) string {
    var output string 
    if val, ok := toolsInfo.Items[name]; ok {
        var magic string
        if val.Magic {
            magic = "is magical"
        } else {
            magic = "is not magical"
        }
        output = fmt.Sprintf("A %s costs %dzm, weighs %d and %s.  It %s.",
        strings.ReplaceAll(name,"-"," "), 
        val.Cost,
        val.Weight,
        magic,
        val.Use)
    }
    return output
}

func getWandMessage(name string) string {
    var output string 
    if val, ok := wandsInfo.Items[name]; ok {
        output = fmt.Sprintf("A %s costs %dzm, weighs %d and has %s starting charges.  It's pattern is %s. "+
        "%s %s", 
        strings.ReplaceAll(name,"-"," "), 
        val.Cost, 
        val.Weight, 
        val.StartingCharges, 
        val.Type,
        val.Effect,
        val.Broken)
    }
    return output
}

func getRingMessage(name string) string {
    var output string 
    if val, ok := ringsInfo.Items[name]; ok {
        output = fmt.Sprintf("A %s costs %dzm and grants %s. %s.",
        strings.ReplaceAll(name,"-"," "), 
        val.Cost, 
        val.ExtrinsicGranted, 
        val.Notes)
    }
    return output
}

func getPropertyMessage(name string) string {
    var output string 
    if val, ok := propsInfo.Items[name]; ok {
        output = fmt.Sprintf("%s %s.  Notable sources include: ",
        strings.Title(strings.ReplaceAll(name,"-"," ")), 
        val.Effect)
        for i,source := range val.Sources {
            output = output + source
            if i != len(val.Sources)-1 {
                output = output + "; "
            } else {
                output = output + "."
            }
            fmt.Println("> " + source)
        }
    }
    return output
}

func getComestibleMessage(name string) string {
    var output string 
    if val, ok := comestiblesInfo.Items[name]; ok {
        var conduct string
        if val.Conduct == "vegan" {
            conduct = ", is vegan,"
        } else if val.Conduct == "vegetarian" {
            conduct = ", is vegetarian"
        }

        output = fmt.Sprintf("A %s costs %dzm, weighs %d, takes %d time to eat%s and grants %d points of nutrition. %s",
        strings.ReplaceAll(name,"-"," "), 
        val.Cost, 
        val.Weight, 
        val.Time,
        conduct,
        val.NutritionalValue,
        val.Effect)
    }
    return output
}

func getPotionMessage(name string) string {
    var output string 
    if val, ok := potionsInfo.Items[name]; ok {
        output = fmt.Sprintf("A potion of %s costs %dzm and weighs %d. %s",
        strings.ReplaceAll(name,"-"," "), 
        val.Cost, 
        val.Weight, 
        val.Effect)
    }
    return output
}

func getMonsterMessage(name string) string {
    var output string 
    if val, ok := monstersInfo.Items[name]; ok {
        var genocidable string
        if val.Genocidable {
            genocidable = "genocidable"
        } else {
            genocidable = "not genocidable"
        }
        var resistances string
        if val.Resistances != "" {
            resistances = " It is resistant to " + val.Resistances + "."
        }
        var resistancesConveyed string
        if val.ResistancesConveyed != "" {
            resistancesConveyed = " It might convey resistance to " + val.ResistancesConveyed + "."
        }
        var corpseSafe string
        if val.CorpseSafe {
            corpseSafe = "safe"
        } else {
            corpseSafe = "not safe"
        }
        var elbereth string
        if val.Elbereth {
            elbereth = "respects"
        } else {
            elbereth = "does not respect"
        }

        output = fmt.Sprintf("A %s has difficulty %d.  It attacks are %s. It " +
        "has speed %d, %d AC, %d MR, weighs %d, has nutritional value %d " +
        "and %s alignment.  It is a %s creature. It is %s.%s%s " + 
        "Its corpse is %s to eat. It %s Elbereth.%s",
        strings.ReplaceAll(name,"-"," "), 
        val.Difficulty, 
        val.Attacks,
        val.Speed,
        val.AC,
        val.MR,
        val.Weight,
        val.NutritionalValue,
        val.Alignment,
        val.Size,
        genocidable,
        resistances,
        resistancesConveyed,
        corpseSafe,
        elbereth,
        val.Extra)
    }
    return output
}

func getWeapons(fname string) *weapons {
    var w *weapons
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &w)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    for k := range w.Items {
        keys = append(keys, k)
    }

    return w
}

func getArmor(fname string) *armors {
    var a *armors
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &a)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return a
}

func getMonsters(fname string) *monsters {
    var m *monsters
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &m)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return m
}

func getTools(fname string) *tools {
    var t *tools
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &t)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return t
}

func getWands(fname string) *wands {
    var w *wands
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &w)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return w
}

func getRings(fname string) *rings {
    var r *rings
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &r)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return r
}

func getProperties(fname string) *properties {
    var p *properties
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &p)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return p
}

func getComestibles(fname string) *comestibles {
    var c *comestibles
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return c
}

func getPotions(fname string) *potions {
    var p *potions
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &p)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return p
}

func getArtifacts(fname string) *artifacts {
    var a *artifacts
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &a)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return a
}

func getAppearances(fname string) *appearances {
    var a *appearances
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &a)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return a
}

func getAllowedChannels(fname string) *allowedChannels {
    var a *allowedChannels
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &a)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return a
}

func parseOracleMessage(c *twitch.Client, message string, user string) {
    if message == "join" {
        for _, allowedChannel := range allowedBroadcasters.Names {
            if user == allowedChannel {
                c.Join(user)
                return
            }
        }
    } else if message == "depart" {
        c.Depart(user)
    }
}

func parseBroadcasterMessage(c *twitch.Client, message string, user string) {
    if message == "oracle-depart" {
        c.Depart(user)
    } else if message == "oracle-update" {
        updateInfo()
    }
}

func parseMessage(c *twitch.Client, m twitch.PrivateMessage) {
    message := m.Message
    channel := m.Channel
    user := m.User.Name
   
    //words := strings.Split(message, " ")
    
    if !strings.HasPrefix(message, "!") {
        return
    } else {
        message = strings.TrimPrefix(message, "!")
    }

    // Deal with requests for the oracle's attention
    if channel == "oracleofdelphibot" {
        parseOracleMessage(c, message, user)
    }

    // Deal with special requests from broadcasters
    if user == channel {
        parseBroadcasterMessage(c, message, user)
    }
    
    // Deal with all other messages
    if _, ok := weaponsInfo.Items[message]; ok {
        c.Say(channel, getWeaponMessage(message))
    } else if _, ok := armorInfo.Items[message]; ok {
        c.Say(channel, getArmorMessage(message))
    } else if _, ok := monstersInfo.Items[message]; ok {
        c.Say(channel, getMonsterMessage(message))
    } else if _, ok := toolsInfo.Items[message]; ok {
        c.Say(channel, getToolMessage(message))
    } else if _, ok := wandsInfo.Items[message]; ok {
        c.Say(channel, getWandMessage(message))
    } else if _, ok := ringsInfo.Items[message]; ok {
        c.Say(channel, getRingMessage(message))
    } else if _, ok := propsInfo.Items[message]; ok {
        c.Say(channel, getPropertyMessage(message))
    } else if _, ok := comestiblesInfo.Items[message]; ok {
        c.Say(channel, getComestibleMessage(message))
    } else if _, ok := potionsInfo.Items[message]; ok {
        c.Say(channel, getPotionMessage(message))
    } else if _, ok := artifactsInfo.Items[message]; ok {
        c.Say(channel, getArtifactMessage(message))
    } else if actualName, ok := appearsAs.Items[message]; ok {
        m.Message = "!"+actualName
        parseMessage(c, m)
    }
}
func updateInfo() {
    // load the information from yaml files containing stats
    allowedBroadcasters = getAllowedChannels("allowed-channels.yaml")
    weaponsInfo = getWeapons("weapons.yaml")
    armorInfo = getArmor("armor.yaml")
    monstersInfo = getMonsters("monsters.yaml")
    toolsInfo = getTools("tools.yaml")
    wandsInfo = getWands("wands.yaml")
    ringsInfo = getRings("rings.yaml")
    propsInfo = getProperties("properties.yaml")
    comestiblesInfo = getComestibles("comestibles.yaml")
    potionsInfo = getPotions("potions.yaml")
    artifactsInfo = getArtifacts("artifacts.yaml")
    appearsAs = getAppearances("appearances.yaml")

    bagSizes := []int{2, 3, 4}
    keyMatching = closestmatch.New(keys, bagSizes)
    fmt.Println(keyMatching.AccuracyMutatingWords())

}

func main() {
    allowedBroadcasters = getAllowedChannels("allowed-channels.yaml")
    // load the information from yaml files containing stats
    weaponsInfo = getWeapons("weapons.yaml")
    armorInfo = getArmor("armor.yaml")
    monstersInfo = getMonsters("monsters.yaml")
    toolsInfo = getTools("tools.yaml")
    wandsInfo = getWands("wands.yaml")
    ringsInfo = getRings("rings.yaml")
    propsInfo = getProperties("properties.yaml")
    comestiblesInfo = getComestibles("comestibles.yaml")
    potionsInfo = getPotions("potions.yaml")
    artifactsInfo = getArtifacts("artifacts.yaml")
    appearsAs = getAppearances("appearances.yaml")

    // find the bot's name, channel's name and oauth from OS env vars
    bot := os.Getenv("TWITCHBOT")
    channel := os.Getenv("TWITCHCHANNEL")
    oauth := os.Getenv("TWITCHOAUTH")

    client := twitch.NewClient(bot, oauth)

    client.OnPrivateMessage(func(message twitch.PrivateMessage) {
        if message.User.Name != bot {
            parseMessage(client, message)
        }
    })

    // If the broadcaster leaves the stream then the bot should too
    client.OnUserPartMessage(func(message twitch.UserPartMessage) {
        if message.User == message.Channel {
            client.Depart(message.Channel)
        }
    })

    client.Join(channel)

    err := client.Connect()
    if err != nil {
        panic(err)
    }
}
