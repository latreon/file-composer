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
}

const formatInfo: Record<string, { name: string; description: string }> = {
    zip: {
        name: 'ZIP',
        description: 'Good balance of compression ratio and compatibility'
    },
    tar: {
        name: 'TAR',
        description: 'Archive without compression, preserves permissions'
    },
    gz: {
        name: 'GZIP',
        description: 'Fast compression, good for text files'
    },
    bz2: {
        name: 'BZIP2',
        description: 'Better compression than GZIP, but slower'
    },
    xz: {
        name: 'XZ',
        description: 'Excellent compression ratio, slower speed'
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

                const formattedFormats = data.formats.map((formatId: string) => ({
                    id: formatId,
                    name: formatInfo[formatId]?.name || formatId.toUpperCase(),
                    description: formatInfo[formatId]?.description || 'Compression format'
                }))

                setFormats(formattedFormats)
            } catch (error) {
                console.error('Error fetching formats:', error)
                // Fallback to hardcoded formats
                setFormats(Object.entries(formatInfo).map(([id, info]) => ({
                    id,
                    name: info.name,
                    description: info.description
                })))
            } finally {
                setLoading(false)
            }
        }

        fetchFormats()
    }, [])

    return (
        <div className="card p-4">
            <h2 className="text-lg font-semibold mb-4">Select Compression Format</h2>

            {loading ? (
                <div className="animate-pulse h-10 bg-gray-200 dark:bg-gray-700 rounded"></div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                    {formats.map((format) => (
                        <button
                            key={format.id}
                            className={`p-3 rounded-lg border-2 transition-all ${selected === format.id
                                    ? 'border-primary bg-primary bg-opacity-10'
                                    : 'border-gray-200 dark:border-gray-700 hover:border-primary'
                                }`}
                            onClick={() => onSelect(format.id)}
                            disabled={disabled}
                        >
                            <div className="flex justify-between items-center">
                                <span className="font-medium">{format.name}</span>
                                {selected === format.id && (
                                    <span className="text-primary">
                                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                            <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                                        </svg>
                                    </span>
                                )}
                            </div>
                            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                                {format.description}
                            </p>
                        </button>
                    ))}
                </div>
            )}
        </div>
    )
} 