package config

type topics map[string]payload

type payload map[string]command

type command []string

func HandleTopics(topics topics) []string {
	topicsSlice := make([]string, 0, len(topics))

	for key := range topics {
		topicsSlice = append(topicsSlice, key)
	}

	return topicsSlice
}

func GetCommand(topics topics, topic string, payload string) []string {
	if val, ok := topics[topic]; ok {
		if val, ok := val[payload]; ok {
			return []string(val)
		}
	}

	return []string{}
}
