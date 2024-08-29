package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Choice represents an item with a weight
type Choice struct {
	Item   string
	Weight int
}

// WeightedRandomPicker handles the weighted random picking
type WeightedRandomPicker struct {
	choices []Choice
	cumSum  []int
	total   int
}

// NewWeightedRandomPicker creates a new WeightedRandomPicker
func NewWeightedRandomPicker(choices []Choice) *WeightedRandomPicker {
	cumSum := make([]int, len(choices))
	total := 0
	for i, choice := range choices {
		total += choice.Weight
		cumSum[i] = total
	}
	return &WeightedRandomPicker{
		choices: choices,
		cumSum:  cumSum,
		total:   total,
	}
}

// Pick randomly picks an item based on weight
func (p *WeightedRandomPicker) Pick() string {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(p.total)
	for i, sum := range p.cumSum {
		if r < sum {
			return p.choices[i].Item
		}
	}
	return ""
}

// GenerateNonRepeatingOffers generates a non-repeating list of offers
func GenerateNonRepeatingOffers(choices []Choice) []string {
	var selectedOffers []string
	remainingOffers := make([]Choice, len(choices))
	copy(remainingOffers, choices)

	for len(remainingOffers) > 0 {
		picker := NewWeightedRandomPicker(remainingOffers)
		selectedOffer := picker.Pick()
		selectedOffers = append(selectedOffers, selectedOffer)

		for i, offer := range remainingOffers {
			if offer.Item == selectedOffer {
				remainingOffers = append(remainingOffers[:i], remainingOffers[i+1:]...)
				break
			}
		}
	}

	return selectedOffers
}

// SimulateClicks simulates user clicks and adjusts weights based on the click data
func SimulateClicks(choices []Choice, targetRatios map[string]float64, iterations int) {
	// Initialize click counters
	clickCounts := map[string]int{
		"A": 0,
		"B": 0,
		"C": 0,
		"D": 0,
	}

	// Simulate the click process
	for i := 0; i < iterations; i++ {
		offerList := GenerateNonRepeatingOffers(choices)
		clickedFirst := offerList[0] // Simulate user clicking the first offer (30% chance)

		// Simulate a 30% click chance for the first offer
		if rand.Float64() < 0.3 {
			clickCounts[clickedFirst]++

			// Simulate a 50% click chance for the second offer if the first one was clicked
			clickedSecond := offerList[1]
			if rand.Float64() < 0.5 {
				clickCounts[clickedSecond]++
			}
		}

		// Adjust weights based on current click data
		totalClicks := clickCounts["A"] + clickCounts["B"] + clickCounts["C"] + clickCounts["D"]
		if totalClicks > 0 {
			for j := range choices {
				actualRatio := float64(clickCounts[choices[j].Item]) / float64(totalClicks)
				targetRatio := targetRatios[choices[j].Item]

				// Adjust the weight to move towards the target ratio
				if actualRatio < targetRatio {
					choices[j].Weight += 1
				} else if actualRatio > targetRatio {
					if choices[j].Weight > 1 { // Ensure weight doesn't go below 1
						choices[j].Weight -= 1
					}
				}
			}
		}
	}

	// Print final click counts and adjusted weights
	fmt.Println("Final click counts:", clickCounts)
	fmt.Println("Final adjusted weights:")
	for _, choice := range choices {
		fmt.Printf("%s: %d\n", choice.Item, choice.Weight)
	}

	// Calculate the final click ratios
	totalClicks := clickCounts["A"] + clickCounts["B"] + clickCounts["C"] + clickCounts["D"]
	fmt.Println("Final click ratios:")
	for _, choice := range choices {
		actualRatio := float64(clickCounts[choice.Item]) / float64(totalClicks)
		fmt.Printf("%s: %.2f%%\n", choice.Item, actualRatio*100)
	}
}

func main() {
	// Define the offers with their respective weights
	choices := []Choice{
		{Item: "A", Weight: 4},
		{Item: "B", Weight: 2},
		{Item: "C", Weight: 2},
		{Item: "D", Weight: 1},
	}

	// Define target ratios based on the initial weights
	totalWeight := 4 + 2 + 2 + 1
	targetRatios := map[string]float64{
		"A": 4.0 / float64(totalWeight), // 4/9 ≈ 0.444
		"B": 2.0 / float64(totalWeight), // 2/9 ≈ 0.222
		"C": 2.0 / float64(totalWeight), // 2/9 ≈ 0.222
		"D": 1.0 / float64(totalWeight), // 1/9 ≈ 0.111
	}

	// Simulate clicks and adjust weights over 100 iterations
	SimulateClicks(choices, targetRatios, 1000)
}
