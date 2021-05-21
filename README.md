<h1>GoCinema - Portal za bioskope baziran na mikroservisnoj arhitekturi</h1>
Predmetni projekat iz Naprednih tehnika programiranja

<h2>Opis aplikacije</h2>
GoCinema je portal koji omogućuje korisnicima da pregledaju, pretražuju bioskopske projekcije i pregledaju, pretražuju i ocenjuju filmove koje nudi bioskop, a i da rezervišu karte za odabrane projekcije. 

Uloge koje postoje u sistemu su registrovani korisnici, prodavci i menadžeri, medjutim portal mogu koristiti i neautentifikovani korisnici.

<h3>Funkcionalnosti neautentifikovanog korisnika</h3>
Neautentifikovani korisnik ima mogućnost da se prijavi i registruje na sistem. Prilikom registracije, neautentifikovani korisnik unosi korisničko ime, email, lozinku, ime i prezime. Nakon uspešne registracije, korisnik može da se uloguje i koristi ostale funkcionalnosti sistema. Neautentifikovani korisnik može da pregleda dostupne filmove i projekcije i da ih pretražuje po različitim kriterijumima.

<h3>Funkcionalnosti registrovanog korisnika</h3>
Registrovani korisnik ima mogućnost da ocenjuje i komentariše filmove. Takodje, registrovani korisnik može i da rezerviše karte. Rezervisanje karte podrazumeva odabir željene projekcije a zatim, i odabir jednog od slobodnih sedišta (za koja karte nisu ni rezervisane ni prodate). Registrovani korisnik ima mogućnost pregleda, pretrage i poništavanja rezervisanih karata pri čemu se rezervacija karte može poništiti najkasnije 24 sata pre početka date projekcije.

<h3>Funkcionalnosti prodavca</h3>
Prodavac može kao i ostali korisnici da pregleda i pretražuje filmove kao i projekcije, takodje može da rezerviše karte za kupce. Nudi mu se i pregled, pretraga i poništavanje rezervisanih i prodatih karata pri čemu se prodata ili rezervisana karta može poništiti najkasnije 24 sata pre početka date projekcije. 
Prodavac može da prodaje karte pri čemu postoje dve vrste prodaje karata. Prva vrsta prodaje je direktna prodaja koja podrazumeva sledeće korake: prodavac odabere konkretnu projekciju, navodi ime i prezime kupca (eventualno korisničkog imena za registrovane korisnike) i odabira jedno od slobodnih sedišta. Druga vrsta prodaje je prodaja rezervisane karte. Prodavac ima i mogućnost da menja karte i izmena karte omogućava promenu sedišta u sali. 

<h3>Funkcionalnosti menadžera</h3>
Menadžer je zadužen za registraciju novih prodavaca i menadžera kao i dodavanje novih filmova i projekcija. Pored toga, menadžer može kreirati sledeće izveštaje: lista prodatih karata za odabran datum prodaje, lista prodatih karata za odabran datum konkretne projekcije, ukupna cena prodatih karata za zadati film u svim projekcijama, ukupan broj i ukupna cena prodatih karata po prodavcima u poslednjih 30 dana itd (inicijalno su navedeni ovi izveštaji medjutim moguća je izmena i dodavanje novih izveštaja).

<h2>Arhitektura sistema</h2>
Arhitektura sistema bi bila zasnovana na mikroservisima za čiji razvoj se koriste jezici Go, Python i Pharo.

<h3>Korisnički servis</h3>
Flask servis koji omogućava registraciju i prijavu korisnika.

<h3>Mikroservis za ocene i komentare</h3>
Go servis koji treba da pokrije ocenjivanje i komentarisanje filmova od strane registrovanih korisnika.

<h3>Mikroservis za analitiku</h3>
Pharo servis koji omogućava prikaz izveštaja koji su definisani u okviru funkcionalnosti menadžera.

<h3>Vue Client</h3>
Omogućava pristup svim funkcionalnostima sistema.

<h3>Cron job-ovi za poništavanje rezervacije</h3>
Flask mikroservis koji bi svakih 15 minuta poništavao rezervacije koje nisu potvrdjene - ako se kupac koji je rezervisao kartu ne pojavi na vreme (pola sata pre početka projekcije), karta se vraća u slobodnu prodaju.

<h3>Mikroservis za rezrevaciju i kupovinu karata</h3>
Go servis koji omogućava registrovanim korisnicima da rezervišu karte a prodavcima da rezervišu i prodaju karte.

<h3>Mikroservis za dodavanje i pretraživanje</h3>
Go servis koji omogućava menadžerima da dodaju nove filmove i projekcije. Ostalim korisnicima u sistemu omogućava pretrage filmova, projekcija i karata.
