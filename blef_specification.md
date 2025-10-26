# BLEF — Book Library Exchange Format

Version **0.2.0** — Spécification

---

## 📌 Objet du document

Le format **BLEF** (Book Library Exchange Format) permet l’échange **interopérable** de bibliothèques personnelles de livres entre plateformes de lecture, outils de gestion et services en ligne.

Ce document décrit :

- le périmètre fonctionnel
- la structure complète du format JSON
- les règles de validation et d’identification
- les cas d’usage couverts
- des exemples d’implémentation

BLEF est conçu pour être **simple, extensible et lisible**, tout en supportant : ✅ livres dans plusieurs listes ✅ notes utilisateur et statut unique ✅ identifiants normalisés (ISBN, UUID) ✅ séries & tomes ✅ gestion facultative des prêts

---

## 🎯 Périmètre & objectifs

BLEF vise à représenter **l’usage lecteur** et non les métadonnées éditoriales exhaustives.

### Inclus

- Bibliothèque personnelle d’un lecteur
- Statut de lecture unique (lus, à lire…)
- Multiples collections / étagères
- Notes, notes privées, tags
- Propriétaire du livre et prêt éventuel

### Hors périmètre (pour futures versions)

- Synchronisation réseau / API
- Gestion d’état de prêt complexe (rappels…)
- Données commerciales ou de recommandation

---

## 🧱 Structure générale du fichier

BLEF est un document JSON contenant **3 sections principales** :

```js
{
  format: "BLEF",
  version: "0.2.0",
  exported_at: datetime,
  user: { ... },
  books: [ ... ],        // métadonnées bibliographiques
  entries: [ ... ],      // données utilisateur par livre
  collections: [ ... ]   // étagères
}
```

---

## 📚 Modèle de données

### 1️⃣ Books — Livres uniques

Contient **une seule fiche** par livre.

Champs clés :

- `id` = identifiant principal (ISBN13 ou UUIDv4)
- `title`, `authors[]`, `language`
- `identifiers` secondaires (ISBN10, OpenLibrary, Wikidata…)
- `edition` (éditeur, année, format, pages)
- `series` facultatif (nom + numéro de tome)

### 2️⃣ Entries — Données utilisateur

Associe **un livre** à **l’expérience personnelle** :

- `status` (enum, un seul statut actif)
- dates de lecture
- note / critique / tags
- propriété et prêt
- appartenance à plusieurs `collections`

### 3️⃣ Collections — Étagères

Définit les différentes listes de l’utilisateur.

Types standardisés (`enum`) :

```
read, reading, to-read, wishlist, owned, custom
```

---

## 🔐 Identification & unicité

| Élément                        | Règle                                   |
| ------------------------------ | --------------------------------------- |
| Livre (`books[].id`)           | ISBN13 **ou** UUIDv4 requis             |
| Un livre dans plusieurs listes | Géré via `entries[].collection_ids[]`   |
| Unicité globale                | `books[].id` **unique** dans le fichier |

---

## ✅ Contrainte statut

Chaque entrée (`entries[]`) doit avoir :

- **exactement un** `status`
- appartenir à ≥ 1 collection

Statuts possibles (`enum`) :

```
read, reading, to-read, abandoned, wishlist
```

---

## 🎧 Séries & sagas

Champ facultatif :

```
series.name
series.volume
```

Permet l’intégration avec des systèmes de suivi de sagas.

---

## 🤝 Gestion des prêts (optionnel)

Simplifiée pour éviter les workflows transactionnels :

```
ownership.owned: boolean
ownership.loaned: {
  status: boolean,
  to?: string,
  date?: date,
  notes?: string
}
```

Ignorable sans casser la validation.

---

## 🧪 JSON Schema officiel

Le JSON Schema complet de BLEF V0.2.0 est fourni **au format JSON** ci‑dessous (copiable tel quel dans un validateur) :

```json
{SCHEMA_PLACEHOLDER}
```

(Le schéma complet sera automatiquement inséré à la prochaine mise à jour du document.)

---

## ✅ Exemple valide minimal

```json
{
  "format": "BLEF",
  "version": "0.2.0",
  "exported_at": "2025-10-26T14:00:00Z",
  "books": [
    {
      "id": "9780156013987",
      "title": "Le Petit Prince",
      "authors": [{ "name": "Antoine de Saint-Exupéry" }],
      "identifiers": { "isbn13": "9780156013987" },
      "language": "fr"
    }
  ],
  "collections": [
    { "id": "lus", "name": "Lus", "type": "read" }
  ],
  "entries": [
    {
      "book_id": "9780156013987",
      "collection_ids": ["lus"],
      "user_data": { "status": "read" }
    }
  ]
}
```

---

## ❌ Exemple invalide

```json
{
  "format": "BLEF",
  "version": "0.2.0",
  "books": []
}
```

Motifs d’invalidité : `exported_at`, `collections`, `entries` manquants ; aucune donnée utilisateur.

---

## 🔄 Compatibilité plateformes

| Plateforme    | Import possible             | Export compatible     | Notes                           |
| ------------- | --------------------------- | --------------------- | ------------------------------- |
| Goodreads     | ✅ via conversion CSV → BLEF | ✅ CSV → BLEF          | Champs utilisateur couverts     |
| Babelio       | ✅ via CSV                   | ✅                     | Utilise EAN → ISBN13 OK         |
| BookWyrm      | 🔶 via API personnalisée    | ✅                     | Alignable avec ActivityPub Book |
| Calibre       | ✅ via plugins               | ✅ ids propres/UUID OK |                                 |
| Inventaire.io | 🔶 via API                  | 🔶                    | IDs Wikidata natifs             |

---

## 🧭 Roadmap

| Version | Objectifs                                        |
| ------- | ------------------------------------------------ |
| 0.2.x   | Stabilisation du modèle + conversions tools      |
| 0.3.0   | Sérialisation OPDS + synchronisation partielle   |
| 1.0.0   | Publication en standard communautaire + API BLEF |

---

## 📄 Licence

Recommandée : **CC0** (usage libre + contributeurs bienvenus)

---

Fin du document — Version 0.2.0

