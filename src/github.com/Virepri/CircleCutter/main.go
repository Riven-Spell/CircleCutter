package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	"sync"
	"sort"
)

var GenerateQueue sync.WaitGroup

type Point struct {
	X,Y float64
}

type Circle struct {
	Point Point
	Radius float64
}

type Solution struct {
	C Circle
	N int
}

type Solutions []Solution

func (S Solutions) Len() int {return len(S)}
func (S Solutions) Swap(i, j int) {S[i], S[j] = S[j], S[i]}
func (S Solutions) Less(i, j int) bool {return S[i].N < S[j].N}

var Quick bool = false

func main(){
	rand.Seed(time.Now().Unix())
	var p int
	fmt.Scan(&p)
	Points := make([]Point,p)
	for k := range Points {
		fmt.Scan(&Points[k].X,&Points[k].Y)
	}

	Population := make([]Circle, 20)
	GenerateQueue.Add(20)
	for k := range Population {
		go GenerateCircle(&Population[k])
	}
	GenerateQueue.Wait()

	o := GeneticAlg(Points, Population, 50000)
	tmp := Circle{}
	if o == tmp {
		fmt.Println("No solution")
	} else {
		fmt.Println(o.Point.X,o.Point.Y)
		fmt.Println(o.Radius)
	}
}

func Distance(p1,p2 Point) float64 {
	return math.Sqrt(math.Pow(math.Abs(p1.X-p2.X),2) + math.Pow(math.Abs(p1.Y-p2.Y),2))
}

func InsideCircle(p Point, c Circle) bool {
	return Distance(c.Point,p) <= c.Radius
}

func InsideBox(c Circle) bool {
	if c.Point.Y + c.Radius > 1 || c.Point.Y - c.Radius < 0 {
		return false
	}
	if c.Point.X + c.Radius > 1 || c.Point.X - c.Radius < 0 {
		return false
	}
	return true
}

func GenerateCircle(c *Circle) {
	o := Circle{
		Point:Point{
			X:rand.Float64(),
			Y:rand.Float64(),
		},
		Radius:rand.Float64(),
	}
	if InsideBox(o) {
		*c = o
		GenerateQueue.Done()
		return
	} else {
		GenerateCircle(c)
	}
}

func GenerateOff(c *Circle, from Circle) {
	o := Circle{
		Point:Point{
			X:from.Point.X,
			Y:from.Point.Y,
		},
		Radius:from.Radius,
	}
	o.Point.Y = ClampF((rand.Float64() - rand.Float64()) + o.Point.X,0,1)
	o.Point.X = ClampF((rand.Float64() - rand.Float64()) + o.Point.Y,0,1)
	o.Radius = ClampF((rand.Float64() - rand.Float64()) + o.Radius,0,1)

	if InsideBox(o) {
		*c = o
		GenerateQueue.Done()
		return
	} else {
		GenerateOff(c,from)
	}
}

func ClampF(n,min,max float64) float64 {
	if n > max {
		return max
	}
	if n < min {
		return min
	}
	return n
}

func CheckSolution(pts []Point, c Circle) (num uint) {
	num = 0
	for _,v := range pts {
		if InsideCircle(v,c) {
			num++
		}
	}
	return num
}

type Circles []Circle

func (C Circles) Len() int { return len(C) }
func (C Circles) Swap(i,j int) {C[i], C[j] = C[j], C[i]}
func (C Circles) Less(i,j int) bool {return C[i].Radius < C[j].Radius}

var Circs Circles = Circles{}

//return nil after so many iterations
func GeneticAlg(pts []Point, pop []Circle, i int) Circle {
	if i == 0 {
		//return Circle{}
		if len(Circs) == 0 {
			return Circle{}
		} else {
			sort.Sort(Circs)
			return Circs[0]
		}
	}

	Sols := make([]int, len(pop))
	for k := range pop {
		Sols[k] = int(CheckSolution(pts, pop[k]))
	}

	if ic := IntsContains(Sols, len(pts)/2); ic != -1 {
		//return pop[ic]
		if !Quick {
			Circs = append(Circs, pop[ic])
			fmt.Println(pop[ic], 50000-i)
		} else {
			fmt.Println("Solved in",50000-i,"genetic cycles")
			return pop[ic]
		}
	}

	//So, no solutions straight up.
	//Express as Solution{} and sort.
	ES := make(Solutions, len(Sols))
	for k,v := range Sols {
		ES[k] = Solution{
			C:pop[k],
			N:v,
		}
	}
	sort.Sort(ES)

	//Generate the seed for the next generation. Make slight changes to these circles and test again.
	Seed := make(Solutions, 4)
	copy(Seed,ES)

	Population := make([]Circle, 20)
	GenerateQueue.Add(20)
	for k := range Population {
		go GenerateOff(&Population[k], Seed[k/5].C)
	}
	GenerateQueue.Wait()

	return GeneticAlg(pts,Population,i-1)
}

func IntsContains(is []int, i int) int {
	for k,v := range is {
		if v == i {
			return k
		}
	}
	return -1
}