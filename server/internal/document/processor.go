package document

import "strings"

func ChunkText(szText string, inMaxChunkSize int) []string {

	szText = strings.TrimSpace(szText)

	if szText == "" {
		return []string{}
	}

	sentences := strings.Split(szText, ". ")

	var chunks []string
	var currentChunk strings.Builder

	for i, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		if i < len(sentences) - 1 && !strings.HasSuffix(sentence, ".") {
			sentence += "."
		}

		if currentChunk.Len() > 0 && currentChunk.Len() + len(sentence) + 1 > inMaxChunkSize {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString(" ")
		}
		currentChunk.WriteString(sentence)
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
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
