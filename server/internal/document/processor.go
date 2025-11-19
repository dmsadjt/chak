package document

import (
	"regexp"
	"strings"
)

func ChunkText(szText string, inMaxChunkSize int) []string {

	szText = strings.TrimSpace(szText)

	if szText == "" {
		return []string{}
	}

	parts := smartSplit(szText)

	var chunks []string
	var currentChunk strings.Builder

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if currentChunk.Len() > 0 && currentChunk.Len()+len(part)+1 > inMaxChunkSize {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n")
		}
		currentChunk.WriteString(part)
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
}

func smartSplit(szText string) []string {
	var partsArray []string

	headerRegex := regexp.MustCompile(`(?m)^#{1,6}\s+.+$`)
	headerIndices := headerRegex.FindAllStringIndex(szText, -1)

	if len(headerIndices) == 0 {
		return splitByParagraphsAndSentences(szText)
	}

	inLastIndex := 0
	for _, indices := range headerIndices {
		if indices[0] > inLastIndex {
			beforeHeader := szText[inLastIndex:indices[0]]
			partsArray = append(partsArray, splitByParagraphsAndSentences(beforeHeader)...)
		}

		header := szText[indices[0]:indices[1]]
		partsArray = append(partsArray, header)

		inLastIndex = indices[1]
	}

	if inLastIndex < len(szText) {
		afterLast := szText[inLastIndex:]
		partsArray = append(partsArray, splitByParagraphsAndSentences(afterLast)...)

	}

	return partsArray
}

func splitByParagraphsAndSentences(szText string) []string {
	szText = strings.TrimSpace(szText)
	if szText == "" {
		return []string{}
	}

	var partsArray []string

	paragraphs := regexp.MustCompile(`\n\n+`).Split(szText, -1)

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}

		if isCodeBlock(paragraph) {
			partsArray = append(partsArray, paragraph)
			continue
		}

		sentences := strings.Split(paragraph, ". ")
		for i, sentence := range sentences {
			sentence = strings.TrimSpace(sentence)
			if sentence == "" {
				continue
			}
			if i < len(sentences) - 1 && !strings.HasSuffix(sentence, ".") {
				sentence += "."
			}

			partsArray = append(partsArray, sentence)
		}
	 }

	 return partsArray
}

func isCodeBlock(szText string) bool {
	if strings.HasPrefix(szText, "```") || strings.HasPrefix(szText, "~~~") {
		return true
	}

	lines := strings.Split(szText, "\n")
	inIndentedLines := 0
	for _, line  := range lines {
		if len(line) > 0 && (strings.HasPrefix(line, "   ") || strings.HasPrefix(line, "\t")) {
			inIndentedLines++
		}
	}

	return inIndentedLines > len(lines)/2
}

func ChunkBySize(szText string, inMaxChunkSize int) []string {
	szText = strings.TrimSpace(szText)

	if len(szText) <= inMaxChunkSize {
		return []string{szText}
	}

	var chunks []string
	for i := 0; i < len(szText); i += inMaxChunkSize {
		end := i + inMaxChunkSize
		if end > len(szText) {
			end = len(szText)
		}
		chunk := szText[i:end]
		chunks = append(chunks, strings.TrimSpace(chunk))
	}

	return chunks
}
