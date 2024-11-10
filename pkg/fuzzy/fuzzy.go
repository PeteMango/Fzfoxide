// fuzzy/fuzzy.go
package fuzzy

func LevenshteinDistance(a, b string) int {
	aLen := len(a)
	bLen := len(b)

	dist := make([][]int, aLen+1)
	for i := range dist {
		dist[i] = make([]int, bLen+1)
	}

	for i := 0; i <= aLen; i++ {
		dist[i][0] = i
	}
	for j := 0; j <= bLen; j++ {
		dist[0][j] = j
	}

	for i := 1; i <= aLen; i++ {
		for j := 1; j <= bLen; j++ {
			if a[i-1] == b[j-1] {
				dist[i][j] = dist[i-1][j-1]
			} else {
				dist[i][j] = min(
					dist[i-1][j]+1,
					dist[i][j-1]+1,
					dist[i-1][j-1]+1,
				)
			}
		}
	}

	return dist[aLen][bLen]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func SimilarityPercentage(a, b string) float64 {
	distance := LevenshteinDistance(a, b)
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	if maxLen == 0 {
		return 100.0
	}
	return (1.0 - float64(distance)/float64(maxLen)) * 100
}
