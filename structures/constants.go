package structures 

//consts for the different paths
const DEFAULT_PATH = "/"
const RENEWABLECURRENT_PATH = "/energy/v1/renewables/current/"
const RENEWABLEHISTORY_PATH = "/energy/v1/renewables/history/"
const NOTIFICATIONS_PATH = "/energy/v1/notifications/"
const STATUS_PATH = "/energy/v1/status/"
const INFO_PATH = "/energy/v1/info/"

//consts for files and URL's
const FILEPATH = "./structures/energyData.csv"
const TESTCOUNTRYFILE = "./countriesData.json"
const COUNTRYSEARCH = "http://129.241.150.113:8080/v3.1/name/"

//consts for sizes
const MAXCACHESIZE = 15
const DAYSTHRESHOLD = 2