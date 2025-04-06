'use client'

import { useState } from 'react'
import FileUpload from '@/components/FileUpload'
import CompressionResult from '@/components/CompressionResult'
import FormatSelector from '@/components/FormatSelector'
import Header from '@/components/Header'
import Footer from '@/components/Footer'

// Response type for the compression API
interface CompressionResponse {
    success: boolean
    message?: string
    downloadLink?: string
    outputSize?: number
    inputSize?: number
}

export default function Home() {
    const [format, setFormat] = useState<string>('')
    const [isLoading, setIsLoading] = useState<boolean>(false)
    const [result, setResult] = useState<CompressionResponse | null>(null)

    const handleCompression = async (file: File) => {
        setIsLoading(true)
        setResult(null)

        try {
            // Create form data to send the file
            const formData = new FormData()
            formData.append('file', file)
            
            // Only append format if it's selected
            if (format) {
                formData.append('format', format)
            }

            // Send file to API for compression
            const response = await fetch('/api/compress', {
                method: 'POST',
                body: formData,
            })

            const data = await response.json()
            setResult(data)
        } catch (error) {
            console.error('Error compressing file:', error)
            setResult({
                success: false,
                message: 'An error occurred while compressing the file',
            })
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <main>
            <Header />

            <div className="max-w-3xl mx-auto px-4">
                {!result && (
                    <>
                        <div className="text-center mb-8">
                            <h2 className="text-3xl md:text-4xl font-bold mb-2">
                                Compress Any File
                            </h2>
                            <p className="text-gray-600 dark:text-gray-300">
                                Drop any file to compress it instantly
                            </p>
                        </div>

                        <FormatSelector
                            selected={format}
                            onSelect={setFormat}
                            disabled={isLoading}
                        />

                        <div className="mt-6">
                            <FileUpload
                                onFileSelect={handleCompression}
                                isLoading={isLoading}
                                format={format}
                            />
                        </div>
                    </>
                )}

                {result && (
                    <div className="mt-8">
                        <CompressionResult
                            result={result}
                            format={format}
                            onReset={() => setResult(null)}
                        />
                    </div>
                )}
            </div>

            <Footer />
        </main>
    )
} 