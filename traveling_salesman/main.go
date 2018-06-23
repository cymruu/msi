package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

var seed = rand.NewSource(time.Now().UnixNano())
var rand1 = rand.New(seed)

type City struct {
	name string
	x    float64
	y    float64
}

func (c *City) distanceTo(d City) float64 {
	return math.Sqrt(math.Pow(c.x-d.x, 2) * math.Pow(c.y-d.y, 2))
}
func (c *City) Rect() image.Rectangle {
	return image.Rect(int(c.x), int(c.y), int(c.x+5), int(c.y+5))
}
func (c *City) Point() image.Point {
	return image.Point{int(c.x), int(c.y)}
}

type Route struct {
	route    []City
	distance float64
	fitness  float64
}

func (r *Route) setOrginCity(start City) {
	r.route = append([]City{start}, r.route...)
}
func (r *Route) setDestinationCity(end City) {
	r.route = append(r.route, end)
}
func (r *Route) calculateDistance() float64 {
	distance := .0
	for i := 1; i < len(r.route); i++ {
		d := r.route[i-1].distanceTo(r.route[i])
		distance += d
	}
	r.distance = distance
	r.fitness = 1 / distance
	return distance
}
func (r *Route) mutable() []City {
	return r.route[1 : len(r.route)-1]
}
func (r *Route) setmutable(route []City) {
	route = append(route, r.route[len(r.route)-1])
	r.route = append([]City{r.route[0]}, route...)
}
func (r *Route) routeToIMG(destintion string, number int) {
	var maxX, maxY = .0, .0
	for _, city := range r.route {
		if city.x > maxX {
			maxX = city.x
		}
		if city.y > maxY {
			maxY = city.y
		}
	}
	img := image.NewRGBA(image.Rect(0, 0, int(maxX+5), int(maxY+5.0)))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 100, 0, 255}}, image.ZP, draw.Src)
	for i, city := range r.route {
		point := city.Rect()
		var pointColor image.Uniform
		if i == 0 {
			pointColor = image.Uniform{color.RGBA{0, 255, 0, 255}}
		} else if i == len(r.route)-1 {
			pointColor = image.Uniform{color.RGBA{0, 0, 0, 255}}
		} else {
			pointColor = image.Uniform{color.RGBA{255, 0, 0, 255}}
		}
		draw.Draw(img, point, &pointColor, image.ZP, draw.Src)
		if i != len(r.route)-1 {
			drawLine(img, city.Point(), r.route[i+1].Point(), color.RGBA{255, 255, 255, 255})
		}
		myfile, _ := os.Create(fmt.Sprintf("./%s/route_%d.png", destintion, number))
		png.Encode(myfile, img)
	}
}
func drawLine(img draw.Image, start, end image.Point,
	fill color.Color) {
	x0, x1 := start.X, end.X
	y0, y1 := start.Y, end.Y
	Δx := math.Abs(float64(x1 - x0))
	Δy := math.Abs(float64(y1 - y0))
	if Δx >= Δy { // shallow slope
		if x0 > x1 {
			x0, y0, x1, y1 = x1, y1, x0, y0
		}
		y := y0
		yStep := 1
		if y0 > y1 {
			yStep = -1
		}
		remainder := float64(int(Δx/2)) - Δx
		for x := x0; x <= x1; x++ {
			img.Set(x, y, fill)
			remainder += Δy
			if remainder >= 0.0 {
				remainder -= Δx
				y += yStep
			}
		}
	} else { // steep slope
		if y0 > y1 {
			x0, y0, x1, y1 = x1, y1, x0, y0
		}
		x := x0
		xStep := 1
		if x0 > x1 {
			xStep = -1
		}
		remainder := float64(int(Δy/2)) - Δy
		for y := y0; y <= y1; y++ {
			img.Set(x, y, fill)
			remainder += Δx
			if remainder >= 0.0 {
				remainder -= Δy
				x += xStep
			}
		}
	}
}

var cities = []City{
	City{name: "Szczecin", x: 5, y: 30},
	City{name: "Gdynia", x: 100, y: 5},
	City{name: "Bydgoszcz", x: 90, y: 70},
	City{name: "Warszawa", x: 170, y: 120},
	City{name: "Łodź", x: 150, y: 140},
	City{name: "Wrocław", x: 60, y: 200},
	City{name: "Lublin", x: 220, y: 190},
	City{name: "Karków", x: 150, y: 290},
	City{name: "Rzeszów", x: 210, y: 290},
}

type populationT []Route

func (p populationT) CalcFitness() {
	for i := 0; i < len(p); i++ {
		p[i].calculateDistance()
	}
}
func CreateRandomRoute(source []City) Route {
	route := Route{}
	route.route = make([]City, len(source))
	perm := rand1.Perm(len(source))
	for i, v := range perm {
		route.route[v] = source[i]
	}
	return route
}
func CreatePopulation(orgin, destination City, populationSize int) populationT {
	population := make(populationT, populationSize)
	for i := 0; i < populationSize; i++ {
		population[i] = CreateRandomRoute(cities)
		population[i].setOrginCity(orgin)
		population[i].setDestinationCity(destination)
	}
	return population
}
func getBest(population populationT) Route {
	bestFitness := .0
	indexOfBest := 0
	for i := 0; i < len(population); i++ {
		if population[i].fitness > bestFitness {
			bestFitness = population[i].fitness
			indexOfBest = i
		}
	}
	return population[indexOfBest]
}
func rouletteWheelSelection(population populationT, maxFitness float64) populationT {
	pool := make(populationT, 0)
	for i := 0; i < len(population); i++ {
		num := int((population[i].fitness / maxFitness) * 100)
		for n := 0; n < num; n++ {
			pool = append(pool, population[i])
		}
	}
	return pool
}
func checkIfInArray(el City, list []City) bool {
	for _, item := range list {
		if el == item {
			return true
		}
	}
	return false
}
func makeChild(p1, p2 Route) Route {
	route := Route{}
	route.setOrginCity(p1.route[0])
	route.setDestinationCity(p1.route[len(p1.route)-1])
	r1 := p1.mutable()
	r2 := p2.mutable()
	fromp1 := rand1.Intn(len(r1))
	newRoute := make([]City, fromp1)
	for i := 0; i < fromp1; i++ {
		newRoute[i] = r1[i]
	}
	for i := 0; i < len(r2); i++ {
		if !checkIfInArray(r2[i], newRoute) {
			newRoute = append(newRoute, r2[i])
		}
	}
	route.setmutable(newRoute)
	return route
}
func naturalSelection(population populationT, desiredPopulationsize int, maxFitness float64) populationT {
	nextGeneration := make(populationT, desiredPopulationsize)
	for i := 0; i < desiredPopulationsize; i++ {
		mom, dad := population[rand.Intn(len(population))], population[rand.Intn(len(population))]
		child := makeChild(mom, dad)
		nextGeneration[i] = child
	}
	nextGeneration.CalcFitness()
	return nextGeneration
}
func main() {
	orgin := City{name: "Katowice", x: 110, y: 300}
	destination := City{name: "Katowice", x: 110, y: 310}
	population := CreatePopulation(orgin, destination, 500)
	population.CalcFitness()
	for i := 0; i < 15; i++ {
		best := getBest(population)
		pool := rouletteWheelSelection(population, best.fitness)
		population = naturalSelection(pool, len(population), best.fitness)
		fmt.Printf("generation: %d distance:%f fitness: %f\n", i, best.distance, best.fitness)
		best.routeToIMG("output/3/", i)
	}
}
