package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

type IOrganism interface {
	CalcFitness() float64
}

type Organism struct {
	DNA     []byte
	Fitness float64
}

func createRandomDNA(dnaLength int) []byte {
	dna := make([]byte, dnaLength)
	for i := 0; i < dnaLength; i++ {
		dna[i] = byte(rand.Intn(95) + 32)
	}
	return dna
}
func createRandomOrganism(dnaLength int) Organism {
	return Organism{
		DNA: createRandomDNA(dnaLength),
	}
}
func (o *Organism) CalcFitness(target []byte) float64 {
	score := 0
	for i := 0; i < len(o.DNA); i++ {
		if o.DNA[i] == target[i] {
			score++
		}
	}
	return float64(score) / float64(len(o.DNA))
}

//mutate create random organism
// func (o *Organism) mutate(MutationRate float64) {
// 	if rand.Float64() < MutationRate {
// 		o.DNA = createRandomDNA(len(o.DNA))
// 	}
// }

//mutate random genes
func (o *Organism) mutate(MutationRate float64) {
	for i := 0; i < len(o.DNA); i++ {
		if rand.Float64() < MutationRate {
			o.DNA[i] = byte(rand.Intn(95) + 32)
		}
	}
}

type populationT []Organism

func CreatePopulation(populationSize int, dnaLength int) populationT {
	population := make(populationT, populationSize)
	for i := 0; i < populationSize; i++ {
		population[i] = createRandomOrganism(dnaLength)
	}
	return population
}
func (p populationT) CalcFitness(target []byte) {
	for i := 0; i < len(p); i++ {
		p[i].Fitness = p[i].CalcFitness(target)
	}
}
func rouletteWheelSelection(population populationT, maxFitness float64) populationT {
	pool := make(populationT, 0)
	for i := 0; i < len(population); i++ {
		num := int((population[i].Fitness / maxFitness) * 100)
		for n := 0; n < num; n++ {
			pool = append(pool, population[i])
		}
	}
	return pool
}
func naturalSelection(population populationT, desiredPopulationsize int, maxFitness float64, target []byte, MutationRate float64) populationT {
	nextGeneration := make(populationT, desiredPopulationsize)
	for i := 0; i < desiredPopulationsize; i++ {
		mom, dad := population[rand.Intn(len(population))], population[rand.Intn(len(population))]
		child := makeChild(mom, dad)
		child.mutate(MutationRate)
		nextGeneration[i] = child
	}
	nextGeneration.CalcFitness(target)
	return nextGeneration
}

func makeChild(p1 Organism, p2 Organism) Organism {
	child := Organism{
		DNA: make([]byte, len(p1.DNA)),
	}
	cut := rand.Intn(len(p1.DNA))
	for i := 0; i < len(p1.DNA); i++ {
		if i > cut {
			child.DNA[i] = p1.DNA[i]
		} else {
			child.DNA[i] = p2.DNA[i]
		}
	}
	return child
}
func getBest(population populationT) Organism {
	bestFitness := .0
	indexOfBest := 0
	for i := 0; i < len(population); i++ {
		if population[i].Fitness > bestFitness {
			bestFitness = population[i].Fitness
			indexOfBest = i
		}
	}
	return population[indexOfBest]
}
func main() {
	startTime := time.Now()
	rand.Seed(time.Now().UnixNano())
	POPULATIONSIZE := 500
	TARGET := []byte("Chcialbym byc marynarzem")
	MutationRate := 0.005
	population := CreatePopulation(POPULATIONSIZE, len(TARGET))
	population.CalcFitness(TARGET)
	found := false
	generation := 0
	for !found {
		generation++
		bestOrganism := getBest(population)
		fmt.Printf("generation %d | best: %s, fitness: %04f\n", generation, string(bestOrganism.DNA), bestOrganism.Fitness)
		if bytes.Compare(bestOrganism.DNA, TARGET) == 0 {
			found = true
		} else {
			bestOrganism = getBest(population)
			pool := rouletteWheelSelection(population, bestOrganism.Fitness)
			population = naturalSelection(pool, len(population), bestOrganism.Fitness, TARGET, MutationRate)
		}
	}
	elapsed := time.Since(startTime)
	fmt.Printf("Time taken %s\n", elapsed)
}
