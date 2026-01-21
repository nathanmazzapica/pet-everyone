import { tool } from "@opencode-ai/plugin"

export default tool({
    description: "Generate a postgresql data source name using python",
    async execute() {
        const result = await Bun.$`python3 .opencode/tools/generate-dsn.py`.text()
        return result.trim()
    },
})
