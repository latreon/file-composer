import { useEffect, useState } from 'react'
import { FiChevronDown } from 'react-icons/fi'

interface FormatSelectorProps {
    selected: string
    onSelect: (format: string) => void
    disabled?: boolean
}

interface Format {
    id: string
    name: string
    description: string
    icon?: string
}

const formatInfo: Record<string, { name: string; description: string; icon?: string }> = {
    '': {
        name: 'Auto Select',
        description: 'Maintain original file format (e.g., PDF stays as PDF)',
        icon: '🤖'
    },
    pdf: {
        name: 'PDF Optimize',
        description: 'Lossless PDF compression while maintaining quality',
        icon: '📄'
    },
    zip: {
        name: 'ZIP',
        description: 'Good balance of compression ratio and compatibility',
        icon: '📦'
    }
}

export default function FormatSelector({ selected, onSelect, disabled = false }: FormatSelectorProps) {
    const [formats, setFormats] = useState<Format[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        // Fetch available formats from the API
        const fetchFormats = async () => {
            try {
                const response = await fetch('/api/formats')
                const data = await response.json()

                // Add the auto-select option
                const allFormats: string[] = ['']
                allFormats.push(...data.formats)

                const formattedFormats = allFormats.map((formatId: string) => ({
                    id: formatId,
                    name: formatInfo[formatId]?.name || formatId.toUpperCase(),
                    description: formatInfo[formatId]?.description || 'Compression format',
                    icon: formatInfo[formatId]?.icon || '📁'
                }))

                setFormats(formattedFormats)
            } catch (error) {
                console.error('Error fetching formats:', error)
                // Fallback to hardcoded formats
                setFormats(Object.entries(formatInfo).map(([id, info]) => ({
                    id,
                    name: info.name,
                    description: info.description,
                    icon: info.icon || '📁'
                })))
            } finally {
                setLoading(false)
            }
        }

        fetchFormats()
    }, [])

    if (loading) {
        return <div className="animate-pulse h-10 bg-gray-200 dark:bg-gray-700 rounded"></div>
    }

    return (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {formats.map((format) => (
                <button
                    key={format.id}
                    onClick={() => !disabled && onSelect(format.id)}
                    className={`p-4 rounded-lg border-2 transition-all ${
                        selected === format.id
                            ? 'border-primary bg-primary bg-opacity-5'
                            : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
                    } ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
                    disabled={disabled}
                >
                    <div className="flex items-center space-x-3">
                        <span className="text-2xl">{format.icon}</span>
                        <div className="text-left">
                            <h3 className="font-medium">{format.name}</h3>
                            <p className="text-sm text-gray-500 dark:text-gray-400">
                                {format.description}
                            </p>
                        </div>
                    </div>
                </button>
            ))}
        </div>
    )
} 