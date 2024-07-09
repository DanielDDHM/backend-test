package models

type Breed struct {
    ID                        int    `json:"id"`
    Species                   string `json:"species"`
    PetSize                   string `json:"pet_size"`
    PetName                   string `json:"pet_name"`
    AverageMaleAdultWeight    int    `json:"average_male_adult_weight"`
    AverageFemaleAdultWeight  int    `json:"average_female_adult_weight"`
}