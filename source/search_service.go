package main

import (
	"math"
	"sort"
	"strings"
)

const IS_DEBUG = true

type search_service struct {
	DocFreq        metric
	InverseDocFreq metric
	TermIDF        map[string]metric
}

type match struct {
	Result string
	Score  float64
}

type doc_tokens map[string][]string
type metric map[string]float64

func GetSearchService(_doc_tokens doc_tokens) *search_service {
	service := search_service{}
	service.InverseDocFreq = CalcIDF(_doc_tokens)
	service.TermIDF = GetTermFreqIDF(_doc_tokens, service.InverseDocFreq)
	return &service
}

func (sv search_service) Search(input_str string, is_asc bool) []string {
	input_tokens := Tokenize(input_str)
	input_tfidf := CalcTF_IDF(input_tokens, sv.InverseDocFreq)
	matches := sv.FindMatches(input_tfidf)
	sorted_matches := SortMatches(matches, is_asc)
	return sorted_matches
}

func (sv search_service) FindMatches(input_tfidf metric) []match {
	matches := []match{}
	for name, tf_idf := range sv.TermIDF {
		cosine_similarity := CalcCosineSimilarity(input_tfidf, tf_idf)
		if cosine_similarity > 0 {
			matches = append(matches, match{
				Result: name,
				Score:  cosine_similarity,
			})
		}
	}
	return matches
}

func SortMatches(matches []match, is_asc bool) []string {
	results := make([]string, 0, len(matches))
	for i := range matches {
		results = append(results, matches[i].Result)
	}
	sort.Slice(results, func(i, j int) bool {
		if is_asc {
			return matches[i].Score > matches[j].Score
		} else {
			return matches[i].Score < matches[j].Score
		}
	})
	return results
}

func GetTermFreqIDF(_doc_tokens doc_tokens, idf metric) map[string]metric {
	tf_idf := map[string]metric{}
	for doc, tokens := range _doc_tokens {
		tf_idf[doc] = CalcTF_IDF(tokens, idf)
	}
	return tf_idf
}

func CalcTF_IDF(tokens []string, idf metric) metric {
	tf := make(metric)
	for _, token := range tokens {
		tf[token] = tf[token] + 1
	}

	tf_idf := make(metric)
	for key := range tf {
		tf[key] = tf[key] / float64(len(tokens))
		tf_idf[key] = tf[key] * idf[key]
	}
	return tf_idf
}

func CalcIDF(_docs_tokens doc_tokens) metric {
	tokens := []string{}
	for _, _tokens := range _docs_tokens {
		for _, token := range _tokens {
			tokens = append(tokens, token)
		}
	}

	df := make(metric)
	for _, token := range tokens {
		if df[token] == 0 {
			df[token] = 1
		}
	}

	idf := make(metric)
	for term, count := range df {
		idf[term] = math.Log(float64(len(_docs_tokens)) / float64(count+1))
	}
	return idf
}

func Standardlize(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	return s
}

func Tokenize(s string) []string {
	delimiters := []string{",", ".", "|"}
	for _, delimiter := range delimiters {
		s = strings.ReplaceAll(s, delimiter, "")
	}
	tokens := strings.Split(s, " ")
	for key := range tokens {
		tokens[key] = Standardlize(tokens[key])
	}
	return tokens
}
func CalcCosineSimilarity(vec1, vec2 map[string]float64) float64 {
	var dot_product, norm1, norm2 float64
	for w, v1 := range vec1 {
		v2 := vec2[w]
		dot_product += v1 * v2
		norm1 += v1 * v1
	}
	for _, v2 := range vec2 {
		norm2 += v2 * v2
	}
	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}
	return dot_product / (math.Sqrt(norm1) * math.Sqrt(norm2))
}
