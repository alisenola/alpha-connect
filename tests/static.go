package tests

var AssetsJSON = []byte(`
{
  "1": {
    "ID": 1,
    "name": "euro",
    "symbol": "EUR"
  },
  "10": {
    "ID": 10,
    "name": "eos",
    "symbol": "EOS"
  },
  "100": {
    "ID": 100,
    "name": "avalanche",
    "symbol": "AVAX"
  },
  "101": {
    "ID": 101,
    "name": "akropolis",
    "symbol": "AKRO"
  },
  "102": {
    "ID": 102,
    "name": "coti",
    "symbol": "COTI"
  },
  "103": {
    "ID": 103,
    "name": "cartesi",
    "symbol": "CTSI"
  },
  "104": {
    "ID": 104,
    "name": "enjin-coin",
    "symbol": "ENJ"
  },
  "105": {
    "ID": 105,
    "name": "terra",
    "symbol": "LUNA"
  },
  "106": {
    "ID": 106,
    "name": "polkadot",
    "symbol": "DOT"
  },
  "107": {
    "ID": 107,
    "name": "wrapped-bitcoin",
    "symbol": "WBTC"
  },
  "108": {
    "ID": 108,
    "name": "theta",
    "symbol": "THETA"
  },
  "109": {
    "ID": 109,
    "name": "status",
    "symbol": "SNT"
  },
  "11": {
    "ID": 11,
    "name": "bitcoin-cash",
    "symbol": "BCH"
  },
  "110": {
    "ID": 110,
    "name": "bancor",
    "symbol": "BNT"
  },
  "111": {
    "ID": 111,
    "name": "neo-gas",
    "symbol": "GAS"
  },
  "112": {
    "ID": 112,
    "name": "oax",
    "symbol": "OAX"
  },
  "113": {
    "ID": 113,
    "name": "district0x",
    "symbol": "DNT"
  },
  "114": {
    "ID": 114,
    "name": "walton-chain",
    "symbol": "WTC"
  },
  "115": {
    "ID": 115,
    "name": "yoyow",
    "symbol": "YOYO"
  },
  "116": {
    "ID": 116,
    "name": "singulardtv",
    "symbol": "SNGLS"
  },
  "117": {
    "ID": 117,
    "name": "myneighboralice",
    "symbol": "ALICE"
  },
  "118": {
    "ID": 118,
    "name": "chiliz",
    "symbol": "CHZ"
  },
  "119": {
    "ID": 119,
    "name": "beefy-finance",
    "symbol": "BIFI"
  },
  "12": {
    "ID": 12,
    "name": "tether",
    "symbol": "USDT"
  },
  "120": {
    "ID": 120,
    "name": "binance-usd",
    "symbol": "BUSD"
  },
  "121": {
    "ID": 121,
    "name": "south-korean-won",
    "symbol": "KRW"
  },
  "122": {
    "ID": 122,
    "name": "1inch",
    "symbol": "1INCH"
  },
  "123": {
    "ID": 123,
    "name": "nkn",
    "symbol": "NKN"
  },
  "124": {
    "ID": 124,
    "name": "origin-protocol",
    "symbol": "OGN"
  },
  "125": {
    "ID": 125,
    "name": "small-love-potion",
    "symbol": "SLP"
  },
  "126": {
    "ID": 126,
    "name": "swipe",
    "symbol": "SXP"
  },
  "127": {
    "ID": 127,
    "name": "kava",
    "symbol": "KAVA"
  },
  "128": {
    "ID": 128,
    "name": "synthetix-network-token",
    "symbol": "SNX"
  },
  "129": {
    "ID": 129,
    "name": "binance-defi-composite-index",
    "symbol": "BINDEFI"
  },
  "13": {
    "ID": 13,
    "name": "stellar",
    "symbol": "XLM"
  },
  "130": {
    "ID": 130,
    "name": "curve-dao-token",
    "symbol": "CRV"
  },
  "131": {
    "ID": 131,
    "name": "tellor",
    "symbol": "TRB"
  },
  "132": {
    "ID": 132,
    "name": "yearn-finance-ii",
    "symbol": "YFII"
  },
  "133": {
    "ID": 133,
    "name": "thorchain",
    "symbol": "RUNE"
  },
  "134": {
    "ID": 134,
    "name": "sushiswap",
    "symbol": "SUSHI"
  },
  "135": {
    "ID": 135,
    "name": "decentraland",
    "symbol": "MANA"
  },
  "136": {
    "ID": 136,
    "name": "serum",
    "symbol": "SRM"
  },
  "137": {
    "ID": 137,
    "name": "harmony",
    "symbol": "ONE"
  },
  "138": {
    "ID": 138,
    "name": "bzx-protocol",
    "symbol": "BZRX"
  },
  "139": {
    "ID": 139,
    "name": "solana",
    "symbol": "SOL"
  },
  "14": {
    "ID": 14,
    "name": "cardano",
    "symbol": "ADA"
  },
  "140": {
    "ID": 140,
    "name": "storj",
    "symbol": "STORJ"
  },
  "141": {
    "ID": 141,
    "name": "chromia",
    "symbol": "CHR"
  },
  "142": {
    "ID": 142,
    "name": "bluzelle",
    "symbol": "BLZ"
  },
  "144": {
    "ID": 144,
    "name": "helium",
    "symbol": "HNT"
  },
  "145": {
    "ID": 145,
    "name": "fantom",
    "symbol": "FTM"
  },
  "146": {
    "ID": 146,
    "name": "flamingo",
    "symbol": "FLM"
  },
  "147": {
    "ID": 147,
    "name": "tomochain",
    "symbol": "TOMO"
  },
  "148": {
    "ID": 148,
    "name": "kusama",
    "symbol": "KSM"
  },
  "149": {
    "ID": 149,
    "name": "near-protocol",
    "symbol": "NEAR"
  },
  "15": {
    "ID": 15,
    "name": "tron",
    "symbol": "TRX"
  },
  "150": {
    "ID": 150,
    "name": "aave",
    "symbol": "AAVE"
  },
  "151": {
    "ID": 151,
    "name": "reserve-rights",
    "symbol": "RSR"
  },
  "152": {
    "ID": 152,
    "name": "polygon",
    "symbol": "MATIC"
  },
  "153": {
    "ID": 153,
    "name": "civic",
    "symbol": "CVC"
  },
  "154": {
    "ID": 154,
    "name": "bella-protocol",
    "symbol": "BEL"
  },
  "155": {
    "ID": 155,
    "name": "certik",
    "symbol": "CTK"
  },
  "156": {
    "ID": 156,
    "name": "axie-infinity",
    "symbol": "AXS"
  },
  "157": {
    "ID": 157,
    "name": "alpha-finance-lab",
    "symbol": "ALPHA"
  },
  "158": {
    "ID": 158,
    "name": "skale-network",
    "symbol": "SKL"
  },
  "159": {
    "ID": 159,
    "name": "the-graph",
    "symbol": "GRT"
  },
  "16": {
    "ID": 16,
    "name": "bitcoin-sv",
    "symbol": "BSV"
  },
  "160": {
    "ID": 160,
    "name": "the-sandbox",
    "symbol": "SAND"
  },
  "161": {
    "ID": 161,
    "name": "ankr",
    "symbol": "ANKR"
  },
  "162": {
    "ID": 162,
    "name": "litentry",
    "symbol": "LIT"
  },
  "163": {
    "ID": 163,
    "name": "unifi-protocol-dao",
    "symbol": "UNFI"
  },
  "164": {
    "ID": 164,
    "name": "dodo",
    "symbol": "DODO"
  },
  "165": {
    "ID": 165,
    "name": "reef",
    "symbol": "REEF"
  },
  "166": {
    "ID": 166,
    "name": "safepal",
    "symbol": "SFP"
  },
  "167": {
    "ID": 167,
    "name": "btc-standard-hashrate-token",
    "symbol": "BTCST"
  },
  "168": {
    "ID": 168,
    "name": "linear",
    "symbol": "LINA"
  },
  "169": {
    "ID": 169,
    "name": "stormx",
    "symbol": "STMX"
  },
  "17": {
    "ID": 17,
    "name": "monero",
    "symbol": "XMR"
  },
  "170": {
    "ID": 170,
    "name": "dent",
    "symbol": "DENT"
  },
  "171": {
    "ID": 171,
    "name": "celer-network",
    "symbol": "CELR"
  },
  "172": {
    "ID": 172,
    "name": "metal",
    "symbol": "MTL"
  },
  "173": {
    "ID": 173,
    "name": "funtoken",
    "symbol": "FUN"
  },
  "174": {
    "ID": 174,
    "name": "sonm",
    "symbol": "SNM"
  },
  "175": {
    "ID": 175,
    "name": "internet-computer",
    "symbol": "ICP"
  },
  "176": {
    "ID": 176,
    "name": "bakery-token",
    "symbol": "BAKE"
  },
  "177": {
    "ID": 177,
    "name": "mask-network",
    "symbol": "MASK"
  },
  "178": {
    "ID": 178,
    "name": "mdex",
    "symbol": "MDX"
  },
  "179": {
    "ID": 179,
    "name": "polkastarter",
    "symbol": "POLS"
  },
  "18": {
    "ID": 18,
    "name": "dash",
    "symbol": "DASH"
  },
  "180": {
    "ID": 180,
    "name": "arweave",
    "symbol": "AR"
  },
  "181": {
    "ID": 181,
    "name": "shiba-inu",
    "symbol": "SHIB"
  },
  "182": {
    "ID": 182,
    "name": "burger-swap",
    "symbol": "BURGER"
  },
  "183": {
    "ID": 183,
    "name": "wazirx",
    "symbol": "WRX"
  },
  "184": {
    "ID": 184,
    "name": "pancake-swap",
    "symbol": "CAKE"
  },
  "185": {
    "ID": 185,
    "name": "ampleforth",
    "symbol": "AMPL"
  },
  "186": {
    "ID": 186,
    "name": "mirror-protocol",
    "symbol": "MIR"
  },
  "187": {
    "ID": 187,
    "name": "fc-barcelona-fan-token",
    "symbol": "BAR"
  },
  "188": {
    "ID": 188,
    "name": "api3",
    "symbol": "API3"
  },
  "189": {
    "ID": 189,
    "name": "keep-network",
    "symbol": "KEEP"
  },
  "19": {
    "ID": 19,
    "name": "iota",
    "symbol": "IOTA"
  },
  "190": {
    "ID": 190,
    "name": "alien-worlds",
    "symbol": "TLM"
  },
  "191": {
    "ID": 191,
    "name": "gitcoin",
    "symbol": "GTC"
  },
  "192": {
    "ID": 192,
    "name": "santiment",
    "symbol": "SAN"
  },
  "193": {
    "ID": 193,
    "name": "streamr",
    "symbol": "DATA"
  },
  "194": {
    "ID": 194,
    "name": "nectar",
    "symbol": "NEC"
  },
  "195": {
    "ID": 195,
    "name": "request",
    "symbol": "REQ"
  },
  "196": {
    "ID": 196,
    "name": "bnktothefuture",
    "symbol": "BFT"
  },
  "197": {
    "ID": 197,
    "name": "odem",
    "symbol": "ODE"
  },
  "198": {
    "ID": 198,
    "name": "aragon",
    "symbol": "ANT"
  },
  "199": {
    "ID": 199,
    "name": "0chain",
    "symbol": "ZCN"
  },
  "2": {
    "ID": 2,
    "name": "pound",
    "symbol": "GBP"
  },
  "20": {
    "ID": 20,
    "name": "neo",
    "symbol": "NEO"
  },
  "200": {
    "ID": 200,
    "name": "utrust",
    "symbol": "UTK"
  },
  "201": {
    "ID": 201,
    "name": "blockv",
    "symbol": "VEE"
  },
  "202": {
    "ID": 202,
    "name": "ors-group",
    "symbol": "ORS"
  },
  "203": {
    "ID": 203,
    "name": "lympo",
    "symbol": "LYM"
  },
  "204": {
    "ID": 204,
    "name": "fetch",
    "symbol": "FET"
  },
  "205": {
    "ID": 205,
    "name": "ftx-token",
    "symbol": "FTT"
  },
  "206": {
    "ID": 206,
    "name": "qash",
    "symbol": "QASH"
  },
  "207": {
    "ID": 207,
    "name": "wax",
    "symbol": "WAXP"
  },
  "208": {
    "ID": 208,
    "name": "essentia",
    "symbol": "ESS"
  },
  "209": {
    "ID": 209,
    "name": "everipedia",
    "symbol": "IQ"
  },
  "21": {
    "ID": 21,
    "name": "ethereum-classic",
    "symbol": "ETC"
  },
  "210": {
    "ID": 210,
    "name": "xriba",
    "symbol": "XRA"
  },
  "211": {
    "ID": 211,
    "name": "parkingo",
    "symbol": "GOT"
  },
  "212": {
    "ID": 212,
    "name": "omni",
    "symbol": "OMNI"
  },
  "213": {
    "ID": 213,
    "name": "kleros",
    "symbol": "PNK"
  },
  "214": {
    "ID": 214,
    "name": "rsk-smart-bitcoin",
    "symbol": "RBTC"
  },
  "215": {
    "ID": 215,
    "name": "tether-eurt",
    "symbol": "EURT"
  },
  "217": {
    "ID": 217,
    "name": "v-systems",
    "symbol": "VSYS"
  },
  "218": {
    "ID": 218,
    "name": "callisto-network",
    "symbol": "CLO"
  },
  "219": {
    "ID": 219,
    "name": "gatetoken",
    "symbol": "GT"
  },
  "22": {
    "ID": 22,
    "name": "maker",
    "symbol": "MKR"
  },
  "220": {
    "ID": 220,
    "name": "bitkan",
    "symbol": "KAN"
  },
  "221": {
    "ID": 221,
    "name": "ultra",
    "symbol": "UOS"
  },
  "222": {
    "ID": 222,
    "name": "renrenbit",
    "symbol": "RRB"
  },
  "223": {
    "ID": 223,
    "name": "just",
    "symbol": "JST"
  },
  "224": {
    "ID": 224,
    "name": "hermez-network",
    "symbol": "HEZ"
  },
  "225": {
    "ID": 225,
    "name": "xinfin-network",
    "symbol": "XDC"
  },
  "226": {
    "ID": 226,
    "name": "pluton",
    "symbol": "PLU"
  },
  "227": {
    "ID": 227,
    "name": "utopia-genesis-foundation",
    "symbol": "UOP"
  },
  "228": {
    "ID": 228,
    "name": "stakenet",
    "symbol": "XSN"
  },
  "229": {
    "ID": 229,
    "name": "celsius",
    "symbol": "CEL"
  },
  "23": {
    "ID": 23,
    "name": "ontology",
    "symbol": "ONT"
  },
  "230": {
    "ID": 230,
    "name": "bridge-mutual",
    "symbol": "BMI"
  },
  "231": {
    "ID": 231,
    "name": "mobilecoin",
    "symbol": "MOB"
  },
  "232": {
    "ID": 232,
    "name": "popsicle-finance",
    "symbol": "ICE"
  },
  "233": {
    "ID": 233,
    "name": "oxygen",
    "symbol": "OXY"
  },
  "234": {
    "ID": 234,
    "name": "everest",
    "symbol": "IDX"
  },
  "235": {
    "ID": 235,
    "name": "quantfury-token",
    "symbol": "QTF"
  },
  "236": {
    "ID": 236,
    "name": "fractal",
    "symbol": "FCL"
  },
  "237": {
    "ID": 237,
    "name": "alchemix",
    "symbol": "ALCX"
  },
  "238": {
    "ID": 238,
    "name": "audius",
    "symbol": "AUDIO"
  },
  "239": {
    "ID": 239,
    "name": "badger-dao",
    "symbol": "BADGER"
  },
  "24": {
    "ID": 24,
    "name": "tezos",
    "symbol": "XTZ"
  },
  "240": {
    "ID": 240,
    "name": "bao-finance",
    "symbol": "BAO"
  },
  "241": {
    "ID": 241,
    "name": "brazilian-digital-token",
    "symbol": "BRZ"
  },
  "242": {
    "ID": 242,
    "name": "coin98",
    "symbol": "C98"
  },
  "243": {
    "ID": 243,
    "name": "clover-finance",
    "symbol": "CLV"
  },
  "244": {
    "ID": 244,
    "name": "convergence",
    "symbol": "CONV"
  },
  "245": {
    "ID": 245,
    "name": "cope",
    "symbol": "COPE"
  },
  "246": {
    "ID": 246,
    "name": "covalent",
    "symbol": "CQT"
  },
  "247": {
    "ID": 247,
    "name": "cream-finance",
    "symbol": "CREAM"
  },
  "248": {
    "ID": 248,
    "name": "dawn-protocol",
    "symbol": "DAWN"
  },
  "249": {
    "ID": 249,
    "name": "dmm-governance",
    "symbol": "DMG"
  },
  "25": {
    "ID": 25,
    "name": "nem",
    "symbol": "XEM"
  },
  "250": {
    "ID": 250,
    "name": "emblem",
    "symbol": "EMB"
  },
  "251": {
    "ID": 251,
    "name": "dydx",
    "symbol": "DYDX"
  },
  "252": {
    "ID": 252,
    "name": "spell-token",
    "symbol": "SPELL"
  },
  "253": {
    "ID": 253,
    "name": "aurory",
    "symbol": "AURY"
  },
  "254": {
    "ID": 254,
    "name": "raydium",
    "symbol": "RAY"
  },
  "255": {
    "ID": 255,
    "name": "illuvium",
    "symbol": "ILV"
  },
  "256": {
    "ID": 256,
    "name": "moeda-loyalty-points",
    "symbol": "MDA"
  },
  "257": {
    "ID": 257,
    "name": "monetha",
    "symbol": "MTH"
  },
  "258": {
    "ID": 258,
    "name": "airswap",
    "symbol": "AST"
  },
  "259": {
    "ID": 259,
    "name": "everex",
    "symbol": "EVX"
  },
  "26": {
    "ID": 26,
    "name": "zcash",
    "symbol": "ZEC"
  },
  "260": {
    "ID": 260,
    "name": "viberate",
    "symbol": "VIB"
  },
  "261": {
    "ID": 261,
    "name": "power-ledger",
    "symbol": "POWR"
  },
  "262": {
    "ID": 262,
    "name": "raiden-network-token",
    "symbol": "RDN"
  },
  "263": {
    "ID": 263,
    "name": "amber",
    "symbol": "AMB"
  },
  "264": {
    "ID": 264,
    "name": "quantstamp",
    "symbol": "QSP"
  },
  "265": {
    "ID": 265,
    "name": "adx-net",
    "symbol": "ADX"
  },
  "266": {
    "ID": 266,
    "name": "cindicator",
    "symbol": "CND"
  },
  "267": {
    "ID": 267,
    "name": "wabi",
    "symbol": "WABI"
  },
  "268": {
    "ID": 268,
    "name": "gifto",
    "symbol": "GTO"
  },
  "269": {
    "ID": 269,
    "name": "aelf",
    "symbol": "ELF"
  },
  "27": {
    "ID": 27,
    "name": "vechain",
    "symbol": "VET"
  },
  "270": {
    "ID": 270,
    "name": "aion",
    "symbol": "AION"
  },
  "271": {
    "ID": 271,
    "name": "solrise-finance",
    "symbol": "SLRS"
  },
  "272": {
    "ID": 272,
    "name": "aavegotchi",
    "symbol": "GHST"
  },
  "273": {
    "ID": 273,
    "name": "mina",
    "symbol": "MINA"
  },
  "274": {
    "ID": 274,
    "name": "looksrare",
    "symbol": "LOOKS"
  },
  "275": {
    "ID": 275,
    "name": "bored-ape-yacht-club",
    "symbol": "BAYC"
  },
  "276": {
    "ID": 276,
    "name": "meebits",
    "symbol": "MEEBTS"
  },
  "277": {
    "ID": 277,
    "name": "decentraland-estate",
    "symbol": "EST"
  },
  "278": {
    "ID": 278,
    "name": "Decentraland",
    "symbol": "LAND"
  },
  "279": {
    "ID": 279,
    "name": "rarible-2",
    "symbol": "RARI-2"
  },
  "28": {
    "ID": 28,
    "name": "crypto-com-chain",
    "symbol": "CRO"
  },
  "280": {
    "ID": 280,
    "name": "rarible-1",
    "symbol": "RARI-1"
  },
  "281": {
    "ID": 281,
    "name": "rarible-3",
    "symbol": "RARI-3"
  },
  "282": {
    "ID": 282,
    "name": "rarible-4",
    "symbol": "RARI-4"
  },
  "283": {
    "ID": 283,
    "name": "rarible-5",
    "symbol": "RARI-5"
  },
  "284": {
    "ID": 284,
    "name": "rarible-6",
    "symbol": "RARI-6"
  },
  "285": {
    "ID": 285,
    "name": "crypto-punks",
    "symbol": "PUNK"
  },
  "286": {
    "ID": 286,
    "name": "sandbox-lands-v2",
    "symbol": "LANDV2"
  },
  "287": {
    "ID": 287,
    "name": "sandbox-lands-v1",
    "symbol": "LANDV1"
  },
  "289": {
    "ID": 289,
    "name": "sandbox-asset",
    "symbol": "SANDASS"
  },
  "29": {
    "ID": 29,
    "name": "basic-attention-token",
    "symbol": "BAT"
  },
  "290": {
    "ID": 290,
    "name": "mutant-ape-yacht-club",
    "symbol": "MAYC"
  },
  "291": {
    "ID": 291,
    "name": "loot",
    "symbol": "LOOT"
  },
  "292": {
    "ID": 292,
    "name": "superrare-v2",
    "symbol": "SUPR-V2"
  },
  "294": {
    "ID": 294,
    "name": "superrare-v1",
    "symbol": "SUPR-V1"
  },
  "295": {
    "ID": 295,
    "name": "clonex",
    "symbol": "CLONEX"
  },
  "296": {
    "ID": 296,
    "name": "Doodles",
    "symbol": "DOODLE"
  },
  "297": {
    "ID": 297,
    "name": "CryptoKitties",
    "symbol": "CKITTY"
  },
  "298": {
    "ID": 298,
    "name": "Cool Cats",
    "symbol": "COOL"
  },
  "299": {
    "ID": 299,
    "name": "DCL Summer 2020",
    "symbol": "DCLDGSMMR2020"
  },
  "3": {
    "ID": 3,
    "name": "canadian-dollar",
    "symbol": "CAD"
  },
  "30": {
    "ID": 30,
    "name": "dogecoin",
    "symbol": "DOGE"
  },
  "300": {
    "ID": 300,
    "name": "decentraland-wearables-dcg",
    "symbol": "DCLDCG"
  },
  "301": {
    "ID": 301,
    "name": "dcl://sugarclub_yumi",
    "symbol": "DCLSGRCLBYM"
  },
  "302": {
    "ID": 302,
    "name": "dcl://digital_alchemy",
    "symbol": "DCLDGTLLCHMY"
  },
  "303": {
    "ID": 303,
    "name": "dcl://ethermon_wearables",
    "symbol": "DCLETHRMNWRBL"
  },
  "304": {
    "ID": 304,
    "name": "decentraland-wearables-halloween-2019",
    "symbol": "DCLHLWN19"
  },
  "305": {
    "ID": 305,
    "name": "dcl://china_flying",
    "symbol": "DCLCHNFLYNG"
  },
  "306": {
    "ID": 306,
    "name": "DCL Dappcraft Moonminer",
    "symbol": "DCLDPPCRFTMNMNR"
  },
  "307": {
    "ID": 307,
    "name": "dcl://dg_fall_2020",
    "symbol": "DCLDGFLL2020"
  },
  "308": {
    "ID": 308,
    "name": "dcl://wz_wonderbot",
    "symbol": "DCLWZWNDRBT"
  },
  "309": {
    "ID": 309,
    "name": "dcl://tech_tribal_marc0matic",
    "symbol": "DCLTCHTRBLMRC0MTC"
  },
  "31": {
    "ID": 31,
    "name": "bitcoin-gold",
    "symbol": "BTG"
  },
  "310": {
    "ID": 310,
    "name": "dcl://cz_mercenary_mtz",
    "symbol": "DCLCZMRCNRYMTZ"
  },
  "311": {
    "ID": 311,
    "name": "dcl://ml_pekingopera",
    "symbol": "DCLMLPKNGPR"
  },
  "312": {
    "ID": 312,
    "name": "dcl://atari_launch",
    "symbol": "DCLTR"
  },
  "313": {
    "ID": 313,
    "name": "dcl://binance_us_collection",
    "symbol": "DCLBNNCSCLLCTN"
  },
  "314": {
    "ID": 314,
    "name": "dcl://xmas_2020",
    "symbol": "DCLXMAS2020"
  },
  "315": {
    "ID": 315,
    "name": "dcl://xmash_up_2020",
    "symbol": "DCLXMSHP2020"
  },
  "316": {
    "ID": 316,
    "name": "dcl://release_the_kraken",
    "symbol": "DCLRLSTHKRKN"
  },
  "317": {
    "ID": 317,
    "name": "dcl://meme_dontbuythis",
    "symbol": "DCLMMDNTBYTHS"
  },
  "318": {
    "ID": 318,
    "name": "dcl://3lau_basics",
    "symbol": "DCL3LBSCS"
  },
  "319": {
    "ID": 319,
    "name": "dcl://rac_basics",
    "symbol": "DCLRCBSCS"
  },
  "32": {
    "ID": 32,
    "name": "omisego",
    "symbol": "OMG"
  },
  "320": {
    "ID": 320,
    "name": "dcl://winklevoss_capital",
    "symbol": "DCLWNKLVSSCPTL"
  },
  "321": {
    "ID": 321,
    "name": "dcl://rtfkt_x_atari",
    "symbol": "DCLRTFTKTI"
  },
  "322": {
    "ID": 322,
    "name": "DCL Moonshot",
    "symbol": "MOONSHOT"
  },
  "323": {
    "ID": 323,
    "name": "decentraland-wearables-stay-safe",
    "symbol": "DCLSS"
  },
  "324": {
    "ID": 324,
    "name": "dcl://halloween_2020",
    "symbol": "DCLHLWN2020"
  },
  "325": {
    "ID": 325,
    "name": "dcl://dg_atari_dillon_francis",
    "symbol": "DCLDGTRDLLNFRCS"
  },
  "326": {
    "ID": 326,
    "name": "dcl://wonderzone_steampunk",
    "symbol": "DCLWNDRZNSTMPNK"
  },
  "327": {
    "ID": 327,
    "name": "decentraland-wearables-xmas-2019",
    "symbol": "DCLXMAS19"
  },
  "328": {
    "ID": 328,
    "name": "dcl://cybermike_cybersoldier_set",
    "symbol": "DCLCYBMKCYBRSLDRST"
  },
  "329": {
    "ID": 329,
    "name": "dcl://dc_meta",
    "symbol": "DCLDCMT"
  },
  "33": {
    "ID": 33,
    "name": "waves",
    "symbol": "WAVES"
  },
  "330": {
    "ID": 330,
    "name": "decentraland-wearables-exclusive-masks",
    "symbol": "DCLEM"
  },
  "331": {
    "ID": 331,
    "name": "DCL Outta This World",
    "symbol": "DCLPMTTTHSWRLD"
  },
  "332": {
    "ID": 332,
    "name": "dcl://dc_niftyblocksmith",
    "symbol": "DCLDCNFTYBLCKSMTH"
  },
  "334": {
    "ID": 334,
    "name": "DCL Wonderzone Meteorchaser",
    "symbol": "DCLWNDRZNMTRCHSR"
  },
  "335": {
    "ID": 335,
    "name": "DCL Dgtble Headspace",
    "symbol": "DCLDGTBLHDSPC"
  },
  "336": {
    "ID": 336,
    "name": "decentraland-wearables-community-contest",
    "symbol": "DCLCC"
  },
  "337": {
    "ID": 337,
    "name": "dcl://pm_dreamverse_eminence",
    "symbol": "DCLPMDRMVRSMNNC"
  },
  "338": {
    "ID": 338,
    "name": "dcl://mf_sammichgamer",
    "symbol": "DCLMFSMMCHGMR"
  },
  "339": {
    "ID": 339,
    "name": "dcl://ml_liondance",
    "symbol": "DCLLMLNDNC"
  },
  "34": {
    "ID": 34,
    "name": "qtum",
    "symbol": "QTUM"
  },
  "340": {
    "ID": 340,
    "name": "decentraland-wearables-launch",
    "symbol": "DCLL"
  },
  "341": {
    "ID": 341,
    "name": "decentraland-wearables-mycryptoheroes",
    "symbol": "DCLMCH"
  },
  "342": {
    "ID": 342,
    "name": "Cryptoadz",
    "symbol": "TOADZ"
  },
  "343": {
    "ID": 343,
    "name": "PUNKS Comic 2",
    "symbol": "COMIC2"
  },
  "344": {
    "ID": 344,
    "name": "PUNKS Comic",
    "symbol": "COMIC"
  },
  "345": {
    "ID": 345,
    "name": "Pixel Vault Founder's DAO",
    "symbol": "PVFD"
  },
  "347": {
    "ID": 347,
    "name": "PUNKS Comic (Special Edition)",
    "symbol": "COMICSE"
  },
  "348": {
    "ID": 348,
    "name": "parallel",
    "symbol": "LL"
  },
  "349": {
    "ID": 349,
    "name": "Sorare New",
    "symbol": "SOR"
  },
  "35": {
    "ID": 35,
    "name": "usd-coin",
    "symbol": "USDC"
  },
  "350": {
    "ID": 350,
    "name": "BoredApeKennelClub",
    "symbol": "BAKC"
  },
  "351": {
    "ID": 351,
    "name": "Azuki",
    "symbol": "AZUKI"
  },
  "352": {
    "ID": 352,
    "name": "World Of Women",
    "symbol": "WOW"
  },
  "353": {
    "ID": 353,
    "name": "MekaVerse",
    "symbol": "MEKA"
  },
  "354": {
    "ID": 354,
    "name": "PudgyPenguins",
    "symbol": "PPG"
  },
  "355": {
    "ID": 355,
    "name": "CyberKongz",
    "symbol": "KONGZ"
  },
  "356": {
    "ID": 356,
    "name": "0N1 Force",
    "symbol": "0N1"
  },
  "357": {
    "ID": 357,
    "name": "Hashmasks",
    "symbol": "HM"
  },
  "358": {
    "ID": 358,
    "name": "VeeFriends",
    "symbol": "VFT"
  },
  "359": {
    "ID": 359,
    "name": "Phanta Bear",
    "symbol": "PHB"
  },
  "36": {
    "ID": 36,
    "name": "decred",
    "symbol": "DCR"
  },
  "360": {
    "ID": 360,
    "name": "my-curio-cards",
    "symbol": "CURIO"
  },
  "361": {
    "ID": 361,
    "name": "Creature World",
    "symbol": "CREATURE"
  },
  "362": {
    "ID": 362,
    "name": "bored-ape-chemistry-club",
    "symbol": "BACC"
  },
  "363": {
    "ID": 363,
    "name": "Emblem Vault V2",
    "symbol": "Emblem.pro"
  },
  "364": {
    "ID": 364,
    "name": "clonex-mintvial",
    "symbol": "CLNXMV"
  },
  "365": {
    "ID": 365,
    "name": "lostpoets",
    "symbol": "LOSTPOETS"
  },
  "366": {
    "ID": 366,
    "name": "Lost Poets",
    "symbol": "POETS"
  },
  "367": {
    "ID": 367,
    "name": "FLUF",
    "symbol": "FLUF"
  },
  "369": {
    "ID": 369,
    "name": "MyCryptoHeroes:LandSector",
    "symbol": "MCHL"
  },
  "37": {
    "ID": 37,
    "name": "lisk",
    "symbol": "LSK"
  },
  "370": {
    "ID": 370,
    "name": "mycryptoheroes",
    "symbol": "MCHE"
  },
  "371": {
    "ID": 371,
    "name": "mycryptoheroes-hero",
    "symbol": "MCHH"
  },
  "373": {
    "ID": 373,
    "name": "Somnium Space Worlds",
    "symbol": "WORLD"
  },
  "374": {
    "ID": 374,
    "name": "somnium-space-cubes",
    "symbol": "SOMNSC"
  },
  "375": {
    "ID": 375,
    "name": "Somnium Space Avatars",
    "symbol": "AVATAR"
  },
  "376": {
    "ID": 376,
    "name": "somnium-space-land",
    "symbol": "SOMNSL"
  },
  "377": {
    "ID": 377,
    "name": "somnium-space-items",
    "symbol": "SOMNSI"
  },
  "379": {
    "ID": 379,
    "name": "gala-games",
    "symbol": "GALAG"
  },
  "38": {
    "ID": 38,
    "name": "trueusd",
    "symbol": "TUSD"
  },
  "380": {
    "ID": 380,
    "name": "Axie",
    "symbol": "AXIE"
  },
  "381": {
    "ID": 381,
    "name": "CyberKongz VX",
    "symbol": "KONGZ VX"
  },
  "382": {
    "ID": 382,
    "name": "Prime Ape Planet",
    "symbol": "PAP"
  },
  "383": {
    "ID": 383,
    "name": "Lazy Lions",
    "symbol": "LION"
  },
  "384": {
    "ID": 384,
    "name": "Cryptovoxels",
    "symbol": "CVPA"
  },
  "385": {
    "ID": 385,
    "name": "HAPE PRIME",
    "symbol": "HAPE"
  },
  "386": {
    "ID": 386,
    "name": "vox-series-1",
    "symbol": "VOX1"
  },
  "387": {
    "ID": 387,
    "name": "SupDucks",
    "symbol": "SD"
  },
  "388": {
    "ID": 388,
    "name": "Sneaky Vampire Syndicate",
    "symbol": "SVS"
  },
  "389": {
    "ID": 389,
    "name": "adidas Originals: Into the Metaverse",
    "symbol": "ADI"
  },
  "39": {
    "ID": 39,
    "name": "ravencoin",
    "symbol": "RVN"
  },
  "390": {
    "ID": 390,
    "name": "KaijuKingz",
    "symbol": "KAIJU"
  },
  "391": {
    "ID": 391,
    "name": "DeadFellaz",
    "symbol": "DEADFELLAZ"
  },
  "392": {
    "ID": 392,
    "name": "Autoglyphs",
    "symbol": "☵"
  },
  "393": {
    "ID": 393,
    "name": "guttercatgang",
    "symbol": "GTRCG"
  },
  "394": {
    "ID": 394,
    "name": "Galactic Apes",
    "symbol": "GALAPE"
  },
  "395": {
    "ID": 395,
    "name": "wolf-game-v2",
    "symbol": "WGAMEV2"
  },
  "397": {
    "ID": 397,
    "name": "Jungle Freaks",
    "symbol": "JFRK"
  },
  "398": {
    "ID": 398,
    "name": "CryptoSkulls",
    "symbol": "CryptoSkulls"
  },
  "399": {
    "ID": 399,
    "name": "TheCurrency",
    "symbol": "TENDER"
  },
  "4": {
    "ID": 4,
    "name": "japenese-yen",
    "symbol": "JPY"
  },
  "40": {
    "ID": 40,
    "name": "nano",
    "symbol": "NANO"
  },
  "400": {
    "ID": 400,
    "name": "MakersPlace V3",
    "symbol": "MKT2"
  },
  "401": {
    "ID": 401,
    "name": "makersplace-erc-721",
    "symbol": "MAKERS"
  },
  "402": {
    "ID": 402,
    "name": "makersplace-v2",
    "symbol": "MKTV2"
  },
  "403": {
    "ID": 403,
    "name": "makersplace",
    "symbol": "MKT"
  },
  "405": {
    "ID": 405,
    "name": "Acclimated​MoonCats",
    "symbol": "😺"
  },
  "406": {
    "ID": 406,
    "name": "MutantCats",
    "symbol": "MUTCATS"
  },
  "407": {
    "ID": 407,
    "name": "vox-series-2",
    "symbol": "VOX2"
  },
  "408": {
    "ID": 408,
    "name": "RumbleKongLeague",
    "symbol": "RKL"
  },
  "409": {
    "ID": 409,
    "name": "Cold Blooded Creepz",
    "symbol": "CBC"
  },
  "41": {
    "ID": 41,
    "name": "augur",
    "symbol": "REP"
  },
  "410": {
    "ID": 410,
    "name": "ALIENFRENS",
    "symbol": "ALIENFRENS"
  },
  "411": {
    "ID": 411,
    "name": "Neo Tokyo: Outer Identities",
    "symbol": "NEOTOI"
  },
  "412": {
    "ID": 412,
    "name": "ZED Horse",
    "symbol": "ZED"
  },
  "413": {
    "ID": 413,
    "name": "Lil Heroes",
    "symbol": "LIL"
  },
  "414": {
    "ID": 414,
    "name": "Capsule",
    "symbol": "CAPSULE"
  },
  "416": {
    "ID": 416,
    "name": "Anonymice",
    "symbol": "MICE"
  },
  "417": {
    "ID": 417,
    "name": "Neo Tokyo: Identities",
    "symbol": "NEOTI"
  },
  "418": {
    "ID": 418,
    "name": "NFT Worlds",
    "symbol": "NFT Worlds"
  },
  "419": {
    "ID": 419,
    "name": "Robotos",
    "symbol": "ROBO"
  },
  "42": {
    "ID": 42,
    "name": "0x",
    "symbol": "ZRX"
  },
  "420": {
    "ID": 420,
    "name": "CryptoBatz by Ozzy Osbourne",
    "symbol": "BATZ"
  },
  "421": {
    "ID": 421,
    "name": "Adam Bomb Squad",
    "symbol": "ABS"
  },
  "422": {
    "ID": 422,
    "name": "wolf-game-v1",
    "symbol": "WGAMEV1"
  },
  "423": {
    "ID": 423,
    "name": "C-01 Official Collection",
    "symbol": "C-01"
  },
  "424": {
    "ID": 424,
    "name": "KIA",
    "symbol": "KIA"
  },
  "425": {
    "ID": 425,
    "name": "Treeverse",
    "symbol": "TRV"
  },
  "426": {
    "ID": 426,
    "name": "The Humanoids ",
    "symbol": "HMNDS"
  },
  "427": {
    "ID": 427,
    "name": "Sevens Token",
    "symbol": "SEVENS"
  },
  "429": {
    "ID": 429,
    "name": "OxyaOriginProject",
    "symbol": "OXYA"
  },
  "43": {
    "ID": 43,
    "name": "zilliqa",
    "symbol": "ZIL"
  },
  "430": {
    "ID": 430,
    "name": "Groupies",
    "symbol": "GROUPIE"
  },
  "432": {
    "ID": 432,
    "name": "Psychedelics Anonymous Genesis",
    "symbol": "PA"
  },
  "433": {
    "ID": 433,
    "name": "Tom Sachs Rocket Components",
    "symbol": "TSRC"
  },
  "434": {
    "ID": 434,
    "name": "Tom Sachs Rockets",
    "symbol": "TSR"
  },
  "435": {
    "ID": 435,
    "name": "The Heart Project",
    "symbol": "HEARTS"
  },
  "436": {
    "ID": 436,
    "name": "Crypto Bull Society",
    "symbol": "CBS"
  },
  "437": {
    "ID": 437,
    "name": "adventure-gold",
    "symbol": "AGLD"
  },
  "438": {
    "ID": 438,
    "name": "bitmax-token",
    "symbol": "ASD"
  },
  "439": {
    "ID": 439,
    "name": "star-atlas",
    "symbol": "ATLAS"
  },
  "44": {
    "ID": 44,
    "name": "icon",
    "symbol": "ICX"
  },
  "440": {
    "ID": 440,
    "name": "bitdao",
    "symbol": "BIT"
  },
  "441": {
    "ID": 441,
    "name": "boba-network",
    "symbol": "BOBA"
  },
  "442": {
    "ID": 442,
    "name": "celo",
    "symbol": "CELO"
  },
  "443": {
    "ID": 443,
    "name": "eden-network",
    "symbol": "EDEN"
  },
  "444": {
    "ID": 444,
    "name": "bonfida",
    "symbol": "FIDA"
  },
  "445": {
    "ID": 445,
    "name": "flow",
    "symbol": "FLOW"
  },
  "446": {
    "ID": 446,
    "name": "gala",
    "symbol": "GALA"
  },
  "447": {
    "ID": 447,
    "name": "holy-trinity",
    "symbol": "HOLY"
  },
  "448": {
    "ID": 448,
    "name": "huobi-token",
    "symbol": "HT"
  },
  "449": {
    "ID": 449,
    "name": "humanscape",
    "symbol": "HUM"
  },
  "45": {
    "ID": 45,
    "name": "bitshares",
    "symbol": "BTS"
  },
  "450": {
    "ID": 450,
    "name": "immutable-x",
    "symbol": "IMX"
  },
  "451": {
    "ID": 451,
    "name": "kin",
    "symbol": "KIN"
  },
  "452": {
    "ID": 452,
    "name": "mercurial-finance",
    "symbol": "MER"
  },
  "453": {
    "ID": 453,
    "name": "mango-markets",
    "symbol": "MNGO"
  },
  "454": {
    "ID": 454,
    "name": "okb",
    "symbol": "OKB"
  },
  "455": {
    "ID": 455,
    "name": "orbs",
    "symbol": "ORBS"
  },
  "456": {
    "ID": 456,
    "name": "pax-gold",
    "symbol": "PAXG"
  },
  "457": {
    "ID": 457,
    "name": "perpetual-protocol",
    "symbol": "PERP"
  },
  "458": {
    "ID": 458,
    "name": "star-atlas-polis",
    "symbol": "POLIS"
  },
  "459": {
    "ID": 459,
    "name": "prometeus",
    "symbol": "PROM"
  },
  "46": {
    "ID": 46,
    "name": "chainlink",
    "symbol": "LINK"
  },
  "460": {
    "ID": 460,
    "name": "pundix",
    "symbol": "PUNDIX"
  },
  "461": {
    "ID": 461,
    "name": "ramp",
    "symbol": "RAMP"
  },
  "462": {
    "ID": 462,
    "name": "render-token",
    "symbol": "RNDR"
  },
  "463": {
    "ID": 463,
    "name": "ronin",
    "symbol": "RON"
  },
  "464": {
    "ID": 464,
    "name": "keeperdao",
    "symbol": "ROOK"
  },
  "465": {
    "ID": 465,
    "name": "oasis-network",
    "symbol": "ROSE"
  },
  "466": {
    "ID": 466,
    "name": "secret",
    "symbol": "SCRT"
  },
  "467": {
    "ID": 467,
    "name": "serum-ecosystem-token",
    "symbol": "SECO"
  },
  "468": {
    "ID": 468,
    "name": "opendao",
    "symbol": "SOS"
  },
  "469": {
    "ID": 469,
    "name": "sirin-labs-token",
    "symbol": "SRN"
  },
  "47": {
    "ID": 47,
    "name": "bytecoin",
    "symbol": "BCN"
  },
  "470": {
    "ID": 470,
    "name": "step-finance",
    "symbol": "STEP"
  },
  "471": {
    "ID": 471,
    "name": "stacks",
    "symbol": "STX"
  },
  "472": {
    "ID": 472,
    "name": "toncoin",
    "symbol": "TONCOIN"
  },
  "473": {
    "ID": 473,
    "name": "truefi-token",
    "symbol": "TRU"
  },
  "474": {
    "ID": 474,
    "name": "bilira",
    "symbol": "TRYB"
  },
  "475": {
    "ID": 475,
    "name": "tulip-protocol",
    "symbol": "TULIP"
  },
  "476": {
    "ID": 476,
    "name": "wrapped-ether",
    "symbol": "WETH"
  },
  "477": {
    "ID": 477,
    "name": "apecoin",
    "symbol": "APE"
  },
  "478": {
    "ID": 478,
    "name": "dusk-network",
    "symbol": "DUSK"
  },
  "479": {
    "ID": 479,
    "name": "iotex",
    "symbol": "IOTX"
  },
  "48": {
    "ID": 48,
    "name": "bitcoin-diamond",
    "symbol": "BCD"
  },
  "480": {
    "ID": 480,
    "name": "automata-network",
    "symbol": "ATA"
  },
  "481": {
    "ID": 481,
    "name": "klaytn",
    "symbol": "KLAY"
  },
  "482": {
    "ID": 482,
    "name": "arpa-chain",
    "symbol": "ARPA"
  },
  "483": {
    "ID": 483,
    "name": "livepeer",
    "symbol": "LPT"
  },
  "484": {
    "ID": 484,
    "name": "anchor-protocol",
    "symbol": "ANC"
  },
  "485": {
    "ID": 485,
    "name": "green-metaverse-token",
    "symbol": "GMT"
  },
  "486": {
    "ID": 486,
    "name": "ethereum-name-service",
    "symbol": "ENS"
  },
  "487": {
    "ID": 487,
    "name": "maps",
    "symbol": "MAPS"
  },
  "488": {
    "ID": 488,
    "name": "serum-locked",
    "symbol": "SRM_LOCKED"
  },
  "489": {
    "ID": 489,
    "name": "wootrade",
    "symbol": "WOO"
  },
  "49": {
    "ID": 49,
    "name": "holo",
    "symbol": "HOT"
  },
  "490": {
    "ID": 490,
    "name": "yield-guild-games",
    "symbol": "YGG"
  },
  "491": {
    "ID": 491,
    "name": "biconomy",
    "symbol": "BICO"
  },
  "492": {
    "ID": 492,
    "name": "constitutiondao",
    "symbol": "PEOPLE"
  },
  "493": {
    "ID": 493,
    "name": "jasmy",
    "symbol": "JASMY"
  },
  "494": {
    "ID": 494,
    "name": "nervos-network",
    "symbol": "CKB"
  },
  "495": {
    "ID": 495,
    "name": "sun-token",
    "symbol": "SUN"
  },
  "496": {
    "ID": 496,
    "name": "rss3",
    "symbol": "RSS3"
  },
  "497": {
    "ID": 497,
    "name": "kadena",
    "symbol": "KDA"
  },
  "498": {
    "ID": 498,
    "name": "biswap",
    "symbol": "BSW"
  },
  "499": {
    "ID": 499,
    "name": "moonbeam",
    "symbol": "GLMR"
  },
  "5": {
    "ID": 5,
    "name": "bitcoin",
    "symbol": "BTC"
  },
  "50": {
    "ID": 50,
    "name": "iostoken",
    "symbol": "IOST"
  },
  "500": {
    "ID": 500,
    "name": "astar",
    "symbol": "ASTR"
  },
  "501": {
    "ID": 501,
    "name": "frax-share",
    "symbol": "FXS"
  },
  "502": {
    "ID": 502,
    "name": "binaryx",
    "symbol": "BNX"
  },
  "503": {
    "ID": 503,
    "name": "alchemy-pay",
    "symbol": "ACH"
  },
  "504": {
    "ID": 504,
    "name": "braintrust",
    "symbol": "BTRST"
  },
  "505": {
    "ID": 505,
    "name": "radicle",
    "symbol": "RAD"
  },
  "506": {
    "ID": 506,
    "name": "deso",
    "symbol": "DESO"
  },
  "507": {
    "ID": 507,
    "name": "pawtocol",
    "symbol": "UPI"
  },
  "508": {
    "ID": 508,
    "name": "rai",
    "symbol": "RAI"
  },
  "509": {
    "ID": 509,
    "name": "nucypher",
    "symbol": "NU"
  },
  "51": {
    "ID": 51,
    "name": "digibyte",
    "symbol": "DGB"
  },
  "510": {
    "ID": 510,
    "name": "liquity",
    "symbol": "LQTY"
  },
  "511": {
    "ID": 511,
    "name": "ribbon-finance",
    "symbol": "RBN"
  },
  "512": {
    "ID": 512,
    "name": "terrausd",
    "symbol": "UST"
  },
  "513": {
    "ID": 513,
    "name": "cryptex-finance",
    "symbol": "CTX"
  },
  "514": {
    "ID": 514,
    "name": "idex",
    "symbol": "IDEX"
  },
  "515": {
    "ID": 515,
    "name": "kryll",
    "symbol": "KRL"
  },
  "516": {
    "ID": 516,
    "name": "harvest-finance",
    "symbol": "FARM"
  },
  "517": {
    "ID": 517,
    "name": "rari-governance-token",
    "symbol": "RGT"
  },
  "518": {
    "ID": 518,
    "name": "barnbridge",
    "symbol": "BOND"
  },
  "519": {
    "ID": 519,
    "name": "maple",
    "symbol": "MPL"
  },
  "52": {
    "ID": 52,
    "name": "aeternity",
    "symbol": "AET"
  },
  "520": {
    "ID": 520,
    "name": "aioz-network",
    "symbol": "AIOZ"
  },
  "521": {
    "ID": 521,
    "name": "propy",
    "symbol": "PRO"
  },
  "522": {
    "ID": 522,
    "name": "quant",
    "symbol": "QNT"
  },
  "523": {
    "ID": 523,
    "name": "polyswarm",
    "symbol": "NCT"
  },
  "524": {
    "ID": 524,
    "name": "ampleforth-governance-token",
    "symbol": "FORTH"
  },
  "525": {
    "ID": 525,
    "name": "rally",
    "symbol": "RLY"
  },
  "526": {
    "ID": 526,
    "name": "origintrail",
    "symbol": "TRAC"
  },
  "527": {
    "ID": 527,
    "name": "fox-token",
    "symbol": "FOX"
  },
  "528": {
    "ID": 528,
    "name": "goldfinch-protocol",
    "symbol": "GFI"
  },
  "529": {
    "ID": 529,
    "name": "dia",
    "symbol": "DIA"
  },
  "53": {
    "ID": 53,
    "name": "verge",
    "symbol": "XVG"
  },
  "530": {
    "ID": 530,
    "name": "lcx",
    "symbol": "LCX"
  },
  "531": {
    "ID": 531,
    "name": "gods-unchained",
    "symbol": "GODS"
  },
  "532": {
    "ID": 532,
    "name": "aergo",
    "symbol": "ARGO"
  },
  "533": {
    "ID": 533,
    "name": "ethernity-chain",
    "symbol": "ERN"
  },
  "534": {
    "ID": 534,
    "name": "wrapped-luna-token",
    "symbol": "WLUNA"
  },
  "535": {
    "ID": 535,
    "name": "wrapped-centrifuge",
    "symbol": "WCFG"
  },
  "536": {
    "ID": 536,
    "name": "inverse-finance",
    "symbol": "INV"
  },
  "537": {
    "ID": 537,
    "name": "bounce-token",
    "symbol": "AUCTION"
  },
  "538": {
    "ID": 538,
    "name": "highstreet",
    "symbol": "HIGH"
  },
  "539": {
    "ID": 539,
    "name": "polymath-network",
    "symbol": "POLY"
  },
  "54": {
    "ID": 54,
    "name": "steem",
    "symbol": "STEEM"
  },
  "540": {
    "ID": 540,
    "name": "superfarm",
    "symbol": "SUPER"
  },
  "541": {
    "ID": 541,
    "name": "tribe",
    "symbol": "TRIBE"
  },
  "542": {
    "ID": 542,
    "name": "shping",
    "symbol": "SHPING"
  },
  "543": {
    "ID": 543,
    "name": "loom-network",
    "symbol": "LOOM"
  },
  "544": {
    "ID": 544,
    "name": "circuits-of-value",
    "symbol": "COVAL"
  },
  "545": {
    "ID": 545,
    "name": "orion-protocol",
    "symbol": "ORN"
  },
  "546": {
    "ID": 546,
    "name": "xyo",
    "symbol": "XYO"
  },
  "547": {
    "ID": 547,
    "name": "aventus",
    "symbol": "AVT"
  },
  "548": {
    "ID": 548,
    "name": "suku",
    "symbol": "SUKU"
  },
  "549": {
    "ID": 549,
    "name": "voyager-token",
    "symbol": "VGX"
  },
  "55": {
    "ID": 55,
    "name": "siacoin",
    "symbol": "SC"
  },
  "550": {
    "ID": 550,
    "name": "measurable-data-token",
    "symbol": "MDT"
  },
  "551": {
    "ID": 551,
    "name": "derivadao",
    "symbol": "DDX"
  },
  "552": {
    "ID": 552,
    "name": "assemble-protocol",
    "symbol": "ASM"
  },
  "553": {
    "ID": 553,
    "name": "rarible",
    "symbol": "RARI"
  },
  "554": {
    "ID": 554,
    "name": "orca",
    "symbol": "ORCA"
  },
  "555": {
    "ID": 555,
    "name": "crypterium",
    "symbol": "CRPT"
  },
  "556": {
    "ID": 556,
    "name": "synapse",
    "symbol": "SYN"
  },
  "557": {
    "ID": 557,
    "name": "moss-carbon-credit",
    "symbol": "MCO2"
  },
  "558": {
    "ID": 558,
    "name": "gyen",
    "symbol": "GYEN"
  },
  "559": {
    "ID": 559,
    "name": "playdapp",
    "symbol": "PLA"
  },
  "56": {
    "ID": 56,
    "name": "paxos-standard-token",
    "symbol": "PAX"
  },
  "560": {
    "ID": 560,
    "name": "quickswap",
    "symbol": "QUICK"
  },
  "561": {
    "ID": 561,
    "name": "function-x",
    "symbol": "FX"
  },
  "562": {
    "ID": 562,
    "name": "amp",
    "symbol": "AMP"
  },
  "563": {
    "ID": 563,
    "name": "metahero",
    "symbol": "HERO"
  },
  "564": {
    "ID": 564,
    "name": "doggy",
    "symbol": "DOGGY"
  },
  "565": {
    "ID": 565,
    "name": "mines-of-dalarnia",
    "symbol": "DAR"
  },
  "566": {
    "ID": 566,
    "name": "convex-finance",
    "symbol": "CVX"
  },
  "567": {
    "ID": 567,
    "name": "dragonchain",
    "symbol": "DRGN"
  },
  "568": {
    "ID": 568,
    "name": "mcdex",
    "symbol": "MCB"
  },
  "569": {
    "ID": 569,
    "name": "meta",
    "symbol": "MTA"
  },
  "57": {
    "ID": 57,
    "name": "dai",
    "symbol": "DAI"
  },
  "58": {
    "ID": 58,
    "name": "loopring",
    "symbol": "LRC"
  },
  "59": {
    "ID": 59,
    "name": "nuls",
    "symbol": "NULS"
  },
  "6": {
    "ID": 6,
    "name": "litecoin",
    "symbol": "LTC"
  },
  "60": {
    "ID": 60,
    "name": "metaverse",
    "symbol": "ETP"
  },
  "600": {
    "ID": 600,
    "name": "dollar",
    "symbol": "USD"
  },
  "602": {
    "ID": 602,
    "name": "project-galaxy",
    "symbol": "GAL"
  },
  "603": {
    "ID": 603,
    "name": "chain",
    "symbol": "XCN"
  },
  "604": {
    "ID": 604,
    "name": "stepn",
    "symbol": "GST"
  },
  "605": {
    "ID": 605,
    "name": "step-app",
    "symbol": "FITFI"
  },
  "606": {
    "ID": 606,
    "name": "creditcoin",
    "symbol": "CTC"
  },
  "607": {
    "ID": 607,
    "name": "3x-short-cardano-token",
    "symbol": "ADABEAR"
  },
  "608": {
    "ID": 608,
    "name": "3x-long-cardano-token",
    "symbol": "ADABULL"
  },
  "609": {
    "ID": 609,
    "name": "neblio",
    "symbol": "NEBL"
  },
  "61": {
    "ID": 61,
    "name": "gnosis",
    "symbol": "GNO"
  },
  "610": {
    "ID": 610,
    "name": "beam",
    "symbol": "BEAM"
  },
  "611": {
    "ID": 611,
    "name": "dock",
    "symbol": "DOCK"
  },
  "613": {
    "ID": 613,
    "name": "gochain",
    "symbol": "GO"
  },
  "614": {
    "ID": 614,
    "name": "nexo",
    "symbol": "NEXO"
  },
  "615": {
    "ID": 615,
    "name": "rei-network",
    "symbol": "REI"
  },
  "616": {
    "ID": 616,
    "name": "lido-dao",
    "symbol": "LDO"
  },
  "619": {
    "ID": 619,
    "name": "trust-wallet-token",
    "symbol": "TWT"
  },
  "62": {
    "ID": 62,
    "name": "melon",
    "symbol": "MLN"
  },
  "620": {
    "ID": 620,
    "name": "alpine-f1-team-fan-token",
    "symbol": "ALPINE"
  },
  "622": {
    "ID": 622,
    "name": "nanobyte-token",
    "symbol": "NBT"
  },
  "625": {
    "ID": 625,
    "name": "injective-protocol",
    "symbol": "INJ"
  },
  "626": {
    "ID": 626,
    "name": "beta-finance",
    "symbol": "BETA"
  },
  "627": {
    "ID": 627,
    "name": "league-of-kingdoms",
    "symbol": "LOKA"
  },
  "628": {
    "ID": 628,
    "name": "ooki-protocol",
    "symbol": "OOKI"
  },
  "63": {
    "ID": 63,
    "name": "cosmos",
    "symbol": "ATOM"
  },
  "631": {
    "ID": 631,
    "name": "voxies",
    "symbol": "VOXEL"
  },
  "632": {
    "ID": 632,
    "name": "joe",
    "symbol": "JOE"
  },
  "634": {
    "ID": 634,
    "name": "cocos-bcx",
    "symbol": "COCOS"
  },
  "635": {
    "ID": 635,
    "name": "cortex",
    "symbol": "CTXC"
  },
  "637": {
    "ID": 637,
    "name": "flux",
    "symbol": "FLUX"
  },
  "639": {
    "ID": 639,
    "name": "3x-short-balancer-token",
    "symbol": "BALBEAR"
  },
  "64": {
    "ID": 64,
    "name": "eidoo",
    "symbol": "EDO"
  },
  "640": {
    "ID": 640,
    "name": "3x-long-balancer-token",
    "symbol": "BALBULL"
  },
  "641": {
    "ID": 641,
    "name": "0-5x-long-balancer-token",
    "symbol": "BALHALF"
  },
  "642": {
    "ID": 642,
    "name": "1x-short-balancer-token",
    "symbol": "BALHEDGE"
  },
  "643": {
    "ID": 643,
    "name": "index-defi",
    "symbol": "DEFI"
  },
  "644": {
    "ID": 644,
    "name": "kilo-shiba-inu",
    "symbol": "KSHIB"
  },
  "645": {
    "ID": 645,
    "name": "index-bitcoin-dominance",
    "symbol": "BTCDOM"
  },
  "646": {
    "ID": 646,
    "name": "ecash",
    "symbol": "XEC"
  },
  "647": {
    "ID": 647,
    "name": "kilo-ecash",
    "symbol": "KXEC"
  },
  "648": {
    "ID": 648,
    "name": "bread",
    "symbol": "BRD"
  },
  "649": {
    "ID": 649,
    "name": "nav-coin",
    "symbol": "NAV"
  },
  "65": {
    "ID": 65,
    "name": "unus-sed-leo",
    "symbol": "LEO"
  },
  "650": {
    "ID": 650,
    "name": "pivx",
    "symbol": "PIVX"
  },
  "651": {
    "ID": 651,
    "name": "wanchain",
    "symbol": "WAN"
  },
  "652": {
    "ID": 652,
    "name": "syscoin",
    "symbol": "SYS"
  },
  "653": {
    "ID": 653,
    "name": "groestlcoin",
    "symbol": "GRS"
  },
  "654": {
    "ID": 654,
    "name": "quarkchain",
    "symbol": "QKC"
  },
  "655": {
    "ID": 655,
    "name": "nexus",
    "symbol": "NXS"
  },
  "656": {
    "ID": 656,
    "name": "selfkey",
    "symbol": "KEY"
  },
  "657": {
    "ID": 657,
    "name": "nebulas-token",
    "symbol": "NAS"
  },
  "658": {
    "ID": 658,
    "name": "fc-porto-fan-token",
    "symbol": "PORTO"
  },
  "659": {
    "ID": 659,
    "name": "santos-fc-fan-token",
    "symbol": "SANTOS"
  },
  "66": {
    "ID": 66,
    "name": "grin",
    "symbol": "GRIN"
  },
  "660": {
    "ID": 660,
    "name": "acala",
    "symbol": "ACA"
  },
  "661": {
    "ID": 661,
    "name": "media-network",
    "symbol": "MEDIA"
  },
  "662": {
    "ID": 662,
    "name": "tilray",
    "symbol": "TLRY"
  },
  "663": {
    "ID": 663,
    "name": "tesla",
    "symbol": "TSLA"
  },
  "664": {
    "ID": 664,
    "name": "twitter",
    "symbol": "TWTR"
  },
  "665": {
    "ID": 665,
    "name": "uber",
    "symbol": "UBER"
  },
  "666": {
    "ID": 666,
    "name": "upbots",
    "symbol": "UBXT"
  },
  "667": {
    "ID": 667,
    "name": "umee",
    "symbol": "UMEE"
  },
  "668": {
    "ID": 668,
    "name": "united-states-oil-fund ",
    "symbol": "USO"
  },
  "669": {
    "ID": 669,
    "name": "moderna",
    "symbol": "MRNA"
  },
  "67": {
    "ID": 67,
    "name": "golem",
    "symbol": "GLM"
  },
  "670": {
    "ID": 670,
    "name": "marinade-staked-sol",
    "symbol": "MSOL"
  },
  "671": {
    "ID": 671,
    "name": "1000-bittorrent",
    "symbol": "KBTT"
  },
  "672": {
    "ID": 672,
    "name": "binance-wrapped-dot",
    "symbol": "BDOT"
  },
  "673": {
    "ID": 673,
    "name": "nano",
    "symbol": "XNO"
  },
  "674": {
    "ID": 674,
    "name": "contentos",
    "symbol": "COS"
  },
  "675": {
    "ID": 675,
    "name": "theta-fuel",
    "symbol": "TFUEL"
  },
  "676": {
    "ID": 676,
    "name": "merit-circle",
    "symbol": "MC"
  },
  "677": {
    "ID": 677,
    "name": "mobox",
    "symbol": "MBOX"
  },
  "678": {
    "ID": 678,
    "name": "prism",
    "symbol": "PRISM"
  },
  "679": {
    "ID": 679,
    "name": "port-finance",
    "symbol": "PORT"
  },
  "68": {
    "ID": 68,
    "name": "hbar",
    "symbol": "HBAR"
  },
  "680": {
    "ID": 680,
    "name": "pfizer",
    "symbol": "PFE"
  },
  "681": {
    "ID": 681,
    "name": "drep",
    "symbol": "DREP"
  },
  "682": {
    "ID": 682,
    "name": "firo",
    "symbol": "FIRO"
  },
  "683": {
    "ID": 683,
    "name": "apple",
    "symbol": "AAPL"
  },
  "684": {
    "ID": 684,
    "name": "amazon",
    "symbol": "AMZN"
  },
  "685": {
    "ID": 685,
    "name": "juventus-fan-token",
    "symbol": "JUV"
  },
  "686": {
    "ID": 686,
    "name": "paris-saint-germain-fan-token",
    "symbol": "PSG"
  },
  "687": {
    "ID": 687,
    "name": "manchester-city-fan-token",
    "symbol": "CITY"
  },
  "688": {
    "ID": 688,
    "name": "realy",
    "symbol": "REAL"
  },
  "689": {
    "ID": 689,
    "name": "hedget",
    "symbol": "HGET"
  },
  "69": {
    "ID": 69,
    "name": "stable-usd",
    "symbol": "USDS"
  },
  "690": {
    "ID": 690,
    "name": "paypal",
    "symbol": "PYPL"
  },
  "691": {
    "ID": 691,
    "name": "solend",
    "symbol": "SLND"
  },
  "692": {
    "ID": 692,
    "name": "numerai-nftees",
    "symbol": "NFTEE"
  },
  "693": {
    "ID": 693,
    "name": "aurora-cannabis-inc ",
    "symbol": "ACB"
  },
  "695": {
    "ID": 695,
    "name": "airbnb",
    "symbol": "ABNB"
  },
  "696": {
    "ID": 696,
    "name": "aleph-im",
    "symbol": "ALEPH"
  },
  "697": {
    "ID": 697,
    "name": "amc-entertainment-holdings",
    "symbol": "AMC"
  },
  "698": {
    "ID": 698,
    "name": "advanced-micro-devices",
    "symbol": "AMD"
  },
  "699": {
    "ID": 699,
    "name": "aphria-inc",
    "symbol": "APHA"
  },
  "7": {
    "ID": 7,
    "name": "ethereum",
    "symbol": "ETH"
  },
  "70": {
    "ID": 70,
    "name": "horizen",
    "symbol": "ZEN"
  },
  "700": {
    "ID": 700,
    "name": "ark-innovation-etf",
    "symbol": "ARKK"
  },
  "701": {
    "ID": 701,
    "name": "alibaba",
    "symbol": "BABA"
  },
  "702": {
    "ID": 702,
    "name": "blackberry",
    "symbol": "BB"
  },
  "703": {
    "ID": 703,
    "name": "bilibili-inc ",
    "symbol": "BILI"
  },
  "704": {
    "ID": 704,
    "name": "bitopro-exchange-token",
    "symbol": "BITO"
  },
  "705": {
    "ID": 705,
    "name": "bitwise-10-crypto-index-fund",
    "symbol": "BITW"
  },
  "706": {
    "ID": 706,
    "name": "blocto",
    "symbol": "BLT"
  },
  "707": {
    "ID": 707,
    "name": "tokenclub",
    "symbol": "TCT"
  },
  "708": {
    "ID": 708,
    "name": "lto-network",
    "symbol": "LTO"
  },
  "709": {
    "ID": 709,
    "name": "moviebloc",
    "symbol": "MBL"
  },
  "71": {
    "ID": 71,
    "name": "komodo",
    "symbol": "KMD"
  },
  "710": {
    "ID": 710,
    "name": "standard-tokenization-protocol",
    "symbol": "STPT"
  },
  "711": {
    "ID": 711,
    "name": "hive-blockchain",
    "symbol": "HIVE"
  },
  "712": {
    "ID": 712,
    "name": "ardor",
    "symbol": "ARDR"
  },
  "713": {
    "ID": 713,
    "name": "pnetwork",
    "symbol": "PNT"
  },
  "714": {
    "ID": 714,
    "name": "vethor-token",
    "symbol": "VTHO"
  },
  "715": {
    "ID": 715,
    "name": "irisnet",
    "symbol": "IRIS"
  },
  "716": {
    "ID": 716,
    "name": "fio-protocol",
    "symbol": "FIO"
  },
  "717": {
    "ID": 717,
    "name": "travala",
    "symbol": "AVA"
  },
  "718": {
    "ID": 718,
    "name": "wrapped-nxm",
    "symbol": "WNXM"
  },
  "719": {
    "ID": 719,
    "name": "moonriver",
    "symbol": "MOVR"
  },
  "72": {
    "ID": 72,
    "name": "vertcoin",
    "symbol": "VTC"
  },
  "720": {
    "ID": 720,
    "name": "keep3rv1",
    "symbol": "KP3R"
  },
  "721": {
    "ID": 721,
    "name": "vulcan-forged-pyr",
    "symbol": "PYR"
  },
  "722": {
    "ID": 722,
    "name": "multichain",
    "symbol": "MULTI"
  },
  "723": {
    "ID": 723,
    "name": "ellipsis-epx",
    "symbol": "EPX"
  },
  "724": {
    "ID": 724,
    "name": "stratis",
    "symbol": "STRAX"
  },
  "725": {
    "ID": 725,
    "name": "optimism-ethereum",
    "symbol": "OP"
  },
  "726": {
    "ID": 726,
    "name": "psy-options",
    "symbol": "PSY"
  },
  "727": {
    "ID": 727,
    "name": "pintu-token",
    "symbol": "PTU"
  },
  "728": {
    "ID": 728,
    "name": "benqi",
    "symbol": "QI"
  },
  "729": {
    "ID": 729,
    "name": "synthetify",
    "symbol": "SNY"
  },
  "73": {
    "ID": 73,
    "name": "spendcoin",
    "symbol": "SPND"
  },
  "730": {
    "ID": 730,
    "name": "ethereum-name-service-nft",
    "symbol": "ENS-NFT"
  },
  "74": {
    "ID": 74,
    "name": "sibcoin",
    "symbol": "SIB"
  },
  "75": {
    "ID": 75,
    "name": "ark",
    "symbol": "ARK"
  },
  "76": {
    "ID": 76,
    "name": "particl",
    "symbol": "PART"
  },
  "77": {
    "ID": 77,
    "name": "six-network",
    "symbol": "SIX"
  },
  "78": {
    "ID": 78,
    "name": "hedge-trade",
    "symbol": "HEDG"
  },
  "79": {
    "ID": 79,
    "name": "stratis",
    "symbol": "STRAT"
  },
  "8": {
    "ID": 8,
    "name": "ripple",
    "symbol": "XRP"
  },
  "80": {
    "ID": 80,
    "name": "reddcoin",
    "symbol": "RDD"
  },
  "81": {
    "ID": 81,
    "name": "dmarket",
    "symbol": "DMT"
  },
  "82": {
    "ID": 82,
    "name": "iexec-rlc",
    "symbol": "RLC"
  },
  "83": {
    "ID": 83,
    "name": "puma-pay",
    "symbol": "PMA"
  },
  "84": {
    "ID": 84,
    "name": "bittorrent",
    "symbol": "BTT"
  },
  "85": {
    "ID": 85,
    "name": "ocean-protocol",
    "symbol": "OCEAN"
  },
  "86": {
    "ID": 86,
    "name": "binance-coin",
    "symbol": "BNB"
  },
  "87": {
    "ID": 87,
    "name": "uniswap",
    "symbol": "UNI"
  },
  "88": {
    "ID": 88,
    "name": "ren",
    "symbol": "REN"
  },
  "89": {
    "ID": 89,
    "name": "yearn-finance",
    "symbol": "YFI"
  },
  "9": {
    "ID": 9,
    "name": "skycoin",
    "symbol": "SKY"
  },
  "90": {
    "ID": 90,
    "name": "orchid",
    "symbol": "OXT"
  },
  "91": {
    "ID": 91,
    "name": "algorand",
    "symbol": "ALGO"
  },
  "92": {
    "ID": 92,
    "name": "kyber-network",
    "symbol": "KNC"
  },
  "93": {
    "ID": 93,
    "name": "compound",
    "symbol": "COMP"
  },
  "94": {
    "ID": 94,
    "name": "band-protocol",
    "symbol": "BAND"
  },
  "95": {
    "ID": 95,
    "name": "numeraire",
    "symbol": "NMR"
  },
  "96": {
    "ID": 96,
    "name": "uma",
    "symbol": "UMA"
  },
  "97": {
    "ID": 97,
    "name": "balancer",
    "symbol": "BAL"
  },
  "98": {
    "ID": 98,
    "name": "filecoin",
    "symbol": "FIL"
  },
  "99": {
    "ID": 99,
    "name": "elrond",
    "symbol": "EGLD"
  }
}
`)

var ChainsJSON = []byte(`
{
  "1": {
    "ID": 1,
    "name": "Ethereum Mainnet",
    "type": "EVM"
  },
  "5": {
    "ID": 5,
    "name": "Starknet Mainnet",
    "type": "SVM"
  },
  "6": {
    "ID": 6,
    "name": "ZKSync",
    "type": "ZKEVM"
  }
}`)

var ProtocolsJSON = []byte(`
{
  "1": {
    "ID": 1,
    "name": "ERC-20"
  },
  "2": {
    "ID": 2,
    "name": "ERC-721"
  },
  "4": {
    "ID": 4,
    "name": "ERC-20"
  },
  "5": {
    "ID": 5,
    "name": "ERC-721"
  },
  "6": {
    "ID": 6,
    "name": "ERC-20"
  },
  "7": {
    "ID": 7,
    "name": "ERC-721"
  }
}
`)
