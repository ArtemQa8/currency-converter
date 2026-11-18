package main

const (
	usd = "USD"
	eur = "EUR"
	rub = "RUB"
	jpy = "JPY"
	krw = "KRW"
	cny = "CNY"
)

var currencyAliases = map[string]string{
	"ДОЛЛАР": usd, "ДОЛЛАРЫ": usd, "ДОЛ": usd,
	"ДОЛАР": usd, "ДОЛЛАРОВ": usd, "ДОЛАРОВ": usd,
	"ДОЛЛ": usd, "БАКС": usd, "БАКСОВ": usd,
	"$": usd, "BUCKS": usd, "DOLLAR": usd, "DOL": usd, "DOLAR": usd,
	"LJK": usd, "LJKK": usd, "ВЩД": usd, "ГЫВ": usd,

	"ЕВРО": eur, "ЕВР": eur,
	"€": eur, "EURO": eur,
	"УГКЩ": eur, "TDHJ": eur, "УГК": eur, "TDH": eur,

	"РУБЛЬ": rub, "РУБЛИ": rub, "РУБ": rub, "РУБЛЕЙ": rub, "РУБЛЯ": rub,
	"₽": rub, "RUBLES": rub,
	"КГИДУЫ": rub, "КГИ": rub, "КУИ": rub, "HE<": rub, "HE,": rub,

	"ЙЕНА": jpy, "ЙЕНЫ": jpy, "ЕНА": jpy, "ЙЕН": jpy, "ЕН": jpy,
	"¥": jpy, "JENA": jpy, "YENA": jpy, "YEN": jpy, "JEN": jpy,
	"ОЗН": jpy, "QTYF": jpy, "НУТФ": jpy, "ОУТФ": jpy,

	"ВОНА": krw, "ВОНЫ": krw, "ВОН": krw,
	"₩": krw, "WON": krw,
	"ЛКЦ": krw, "DJYF": krw, "DJY": krw, "ЦЩТ": krw,

	"ЮАНЬ": cny, "ЮАНЕЙ": cny, "ЮАНЯ": cny, "ЮАНИ": cny,
	"Ұ": cny, "YUAN": cny, "RENMINBI": cny, "RMB": cny,
	"ЖЕНЬМИНЬБИ": cny, "РЕНМИНБИ": cny, "ЖМБ": cny, "РМБ": cny,
	">FYM": cny, "НГФТ": cny,
}
