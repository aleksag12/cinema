<h1>GoCinema - Portal za bioskope baziran na mikroservisnoj arhitekturi</h1>
Predmetni projekat iz Naprednih tehnika programiranja

<h2>Opis aplikacije</h2>
GoCinema je portal koji omogućuje korisnicima da pregledaju, pretražuju bioskopske projekcije i pregledaju, pretražuju i ocenjuju filmove koje nudi bioskop, a i da rezervišu karte za odabrane projekcije. 

Uloge koje postoje u sistemu su registrovani korisnici, prodavci i menadžeri, medjutim portal mogu koristiti i neautentifikovani korisnici.

<h3>Funkcionalnosti neautentifikovanog korisnika</h3>
Neautentifikovani korisnik ima mogućnost da se prijavi i registruje na sistem. Prilikom registracije, neautentifikovani korisnik unosi korisničko ime, email, lozinku, ime i prezime. Nakon uspešne registracije, korisnik može da se uloguje i koristi ostale funkcionalnosti sistema. Neautentifikovani korisnik može da pregleda dostupne filmove i projekcije i da ih pretražuje i sortira po različitim kriterijumima.

<h3>Funkcionalnosti registrovanog korisnika</h3>
Registrovani korisnik ima mogućnost da ocenjuje i komentariše filmove kao i da obriše svoje komentare. Takodje, registrovani korisnik može i da rezerviše karte. Rezervisanje karte podrazumeva odabir željene projekcije a zatim, i odabir jednog od slobodnih sedišta (za koja karte nisu ni rezervisane ni prodate). Registrovani korisnik ima mogućnost pregleda i poništavanja rezervisanih karata pri čemu se rezervacija karte može poništiti najkasnije 24 sata pre početka date projekcije.

<h3>Funkcionalnosti prodavca</h3>
Prodavac može kao i ostali korisnici da pregleda filmove i projekcije, takodje može da rezerviše karte za kupce kao i da prodaje karte. Pored toga, nudi mu se i pregled i poništavanje rezervisanih i prodatih karata pri čemu se prodata ili rezervisana karta može poništiti najkasnije 24 sata pre početka date projekcije. 
Prodavac može da prodaje karte pri čemu postoje dve vrste prodaje karata. Prva vrsta prodaje je direktna prodaja koja podrazumeva sledeće korake: prodavac odabere konkretnu projekciju, navodi ime i prezime kupca (eventualno korisničkog imena za registrovane korisnike) i odabira jedno od slobodnih sedišta (za koja karte nisu ni rezervisane ni prodate). Druga vrsta prodaje je prodaja prethodno rezervisane karte.

<h3>Funkcionalnosti menadžera</h3>
Menadžer je zadužen za registraciju novih prodavaca i menadžera kao i dodavanje novih filmova i projekcija. Takodje, menadžer može da menja i briše filmove, pri čemu film ne može biti obrisan ako za njega postoji projekcija. Menadžer može obrisati i projekciju, ali samo u slučaju kada nema prodatih ili rezervisanih karata za tu projekciju. Pored toga, menadžer može kreirati sledeće izveštaje: pregled filmova sa najviše prodatih karata u poslednjih 7 dana u obliku pite kao i pregled filmova od kojih je ostvaren najveći profit od prodaje karata u obliku bar chart-a.

<h2>Arhitektura sistema</h2>
Arhitektura sistema bi bila zasnovana na mikroservisima za čiji razvoj se koriste jezici Go, Python i Pharo.

<h3>API Gateway</h3>
Flask servis koji omogućava uniformni pristup svim funkcionalnostima sistema.

<h3>Korisnički servis</h3>
Flask servis koji omogućava registraciju i prijavu korisnika kao i dodavanje novih prodavaca i menadžera u sistem.

<h3>Mikroservis za filmove</h3>
Go servis koji omogućava dodavanje, izmenu i brisanje filmova.

<h3>Mikroservis za konkretne bioskopske projekcije</h3>
Go servis koji omogućava dodavanje, izmenu i brisanje konkretnih bioskopskih projekcija.

<h3>Mikroservis za ocene i komentare</h3>
Go servis koji treba da pokrije ocenjivanje i komentarisanje filmova od strane registrovanih korisnika, kao i brisanje komentara.

<h3>Mikroservis za rezervaciju i kupovinu karata</h3>
Go servis koji omogućava registrovanim korisnicima da rezervišu karte i otkazuju rezervacije, a prodavcima da rezervišu, prodaju i otkazuju karte.

<h3>Cron job-ovi za poništavanje rezervacije</h3>
Go mikroservis koji na svakih 30 minuta poništava rezervacije koje nisu potvrdjene - ako se kupac koji je rezervisao kartu ne pojavi na vreme (pola sata pre početka projekcije), karta se vraća u slobodnu prodaju.

<h3>Mikroservis za analitiku</h3>
Pharo servis koji omogućava prikaz izveštaja koji su navedeni u okviru funkcionalnosti menadžera.

<h3>Vue Client</h3>
Omogućava pristup svim funkcionalnostima sistema.
