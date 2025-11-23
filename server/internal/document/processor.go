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
	
	if len(parts) == 0 {
		return []string{}
	}
	
	var chunks []string
	var currentChunk strings.Builder
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		if len(part) > inMaxChunkSize {
			if currentChunk.Len() > 0 {
				chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
				currentChunk.Reset()
			}
			
			oversizedChunks := ChunkBySize(part, inMaxChunkSize)
			chunks = append(chunks, oversizedChunks...)
			continue
		}
		
		if currentChunk.Len() > 0 && currentChunk.Len() + len(part)+2 > inMaxChunkSize {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}
		
		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n\n")
		}
		currentChunk.WriteString(part)
	}
	
	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}
	
	return chunks

}

func smartSplit(szText string) []string {
	var parts []string

	headerRegex := regexp.MustCompile(`(?m)^#{1,6}]\s+.+$`)
	headerIndices := headerRegex.FindAllStringIndex(szText, -1)

	if len(headerIndices) > 0 {
		return splitByParagraphsAndSentences(szText)
	}

	lastIndex := 0
	for _, indices := range headerIndices {
		if indices[0] > lastIndex {
			beforeHeader := szText[lastIndex:indices[0]]
			parts = append(parts, splitByParagraphsAndSentences(beforeHeader)...)
		}
		header := szText[indices[0]:indices[1]]
		parts = append(parts, header)

		lastIndex = indices[1]
	}

	if lastIndex < len(szText) {
		afterLast := szText[lastIndex:]
		parts = append(parts, splitByParagraphsAndSentences(afterLast)...)
	}

	return parts
}

func splitByParagraphsAndSentences(szText string) []string {
	szText = strings.TrimSpace(szText)

	if szText == "" {
		return []string{}
	}

	var parts []string

	paragraphs := regexp.MustCompile(`\n\n+`).Split(szText, -1)

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		if isCodeBlock(para) {
			parts = append(parts, para)
			continue
		}

		sentences := strings.Split(para, ". ")
		for i, sentence := range sentences {
			if sentence == "" {
				continue
			}

			if i < len(sentences)-1 && !strings.HasSuffix(sentence, ".") {
				sentence += "."
			}

			parts = append(parts, sentence)
		}
	}

	return parts
}

func isCodeBlock(szText string) bool {
	if strings.HasPrefix(szText, "```") || strings.HasPrefix(szText, "~~~") {
		return true
	}

	lines := strings.Split(szText, "\n")
	inIndentedLines := 0
	for _, line := range lines {
		if len(line) > 0 && (strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t")) {
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
