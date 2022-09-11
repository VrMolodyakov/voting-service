## Vote service


[![Build Status](https://scrutinizer-ci.com/g/VrMolodyakov/voting-service/badges/build.png?b=master&s=f29b6e94360ec9cbeaa37e50f88c7fb9164e7228)](https://scrutinizer-ci.com/g/VrMolodyakov/voting-service/build-status/master)

## Installation


```sh
git clone https://github.com/VrMolodyakov/voting-service.git
make start
```

## Unit test

```sh
make test 
```

| package | Coverage |
| ------ | ------ |
| choiceCache  |  80.0% of statements |
| psqlStorage | 92.0% of statements |
| service | 98.6% of statements |
| handler  | 87.8% of statements |


## Api

```
Post /api/vote
```
Creates a new vote pool
Request Body :
 - vote - title of vote
 - choices - voting options

Example :

```
{
    "vote":  "Best pokemon",
    "choices": [
        "Pikachu",
        "Mew",
        "Mewtwo"
    ]
}
```
Response :

```
{
   "vote": "Best pokemon",
   "choices": [
      {
         "choice": "Pikachu",
         "vote_count": 0
      },
      {
         "choice": "Mew",
         "vote_count": 0
      },
      {
         "choice": "Mewtwo",
         "vote_count": 0
      }
   ]
}
```

```
Post /api/result
```
Get vote result.
Request body:
 - vote - vote title

Example:

```
{
    "vote":"Best pokemon"
}
```
Response:

```
[
   {
      "choice": "Pikachu",
      "vote_count": 0
   },
   {
      "choice": "Mew",
      "vote_count": 0
   },
   {
      "choice": "Mewtwo",
      "vote_count": 0
   }
]
```
To make vote
```
Post /api/choice
```
Request body:
 - vote - vote title
 - choice - choice title

Example :

```
{"vote":"Best pokemon","choice":"Pikachu"}
```

Response :

```
204 status
```

