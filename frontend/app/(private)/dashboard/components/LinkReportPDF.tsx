"use client"

import { Link } from "@/api/links"
import { Button } from "@/components/ui/button"
import { Download } from "lucide-react"
import { useState } from "react"
import { Document, Page, Text, View, StyleSheet, PDFViewer } from '@react-pdf/renderer'

interface LinkReportPDFProps {
    links: Link[]
}

const styles = StyleSheet.create({
    page: {
        padding: 30,
    },
    title: {
        fontSize: 24,
        marginBottom: 20,
    },
    summary: {
        marginBottom: 20,
        fontSize: 12,
    },
    table: {
        display: "flex",
        flexDirection: "column",
        width: "100%",
    },
    tableRow: {
        flexDirection: "row",
        borderBottomWidth: 1,
        borderBottomColor: "#E5E7EB",
        paddingVertical: 8,
    },
    tableHeader: {
        backgroundColor: "#F3F4F6",
        fontWeight: "bold",
    },
    tableCell: {
        flex: 1,
        fontSize: 12,
        paddingHorizontal: 5,
    },
    expiredText: {
        color: "#DC2626",
    },
    activeText: {
        color: "#000000",
    },
})

function LinkReportDocument({ links }: LinkReportPDFProps) {
    const totalLinks = links.length
    const activeLinks = links.filter(link => !link.expiration_date || new Date(link.expiration_date) >= new Date()).length
    const totalClicks = links.reduce((sum, link) => sum + (link.clicks || 0), 0)
    const expiredLinks = links.filter(link => link.expiration_date && new Date(link.expiration_date) < new Date()).length

    return (
        <Document>
            <Page size="A4" style={styles.page}>
                <Text style={styles.title}>Link Performance Report</Text>

                <View style={styles.summary}>
                    <Text>Total Links: {totalLinks}</Text>
                    <Text>Active Links: {activeLinks}</Text>
                    <Text>Expired Links: {expiredLinks}</Text>
                    <Text>Total Clicks: {totalClicks}</Text>
                </View>

                <View style={styles.table}>
                    <View style={[styles.tableRow, styles.tableHeader]}>
                        <Text style={styles.tableCell}>Short URL</Text>
                        <Text style={styles.tableCell}>Clicks</Text>
                        <Text style={styles.tableCell}>Created At</Text>
                        <Text style={styles.tableCell}>Expiration Date</Text>
                        <Text style={styles.tableCell}>Status</Text>
                    </View>
                    {links.map((link) => {
                        const isExpired = link.expiration_date && new Date(link.expiration_date) < new Date()
                        return (
                            <View key={link.id} style={styles.tableRow}>
                                <Text style={styles.tableCell}>{link.short_url}</Text>
                                <Text style={styles.tableCell}>{link.clicks}</Text>
                                <Text style={styles.tableCell}>
                                    {new Date(link.created_at).toLocaleDateString()}
                                </Text>
                                <Text style={styles.tableCell}>
                                    {link.expiration_date
                                        ? new Date(link.expiration_date).toLocaleDateString()
                                        : 'Never'}
                                </Text>
                                <Text style={isExpired ? styles.expiredText : styles.activeText}>
                                    {isExpired ? 'Expired' : 'Active'}
                                </Text>
                            </View>
                        )
                    })}
                </View>
            </Page>
        </Document>
    )
}

export function LinkReportPDF({ links }: LinkReportPDFProps) {
    const [showPDF, setShowPDF] = useState(false)

    return (
        <>
            <Button
                variant="outline"
                onClick={() => setShowPDF(true)}
                className="flex items-center gap-2"
            >
                <Download className="h-4 w-4" />
                Export PDF Report
            </Button>

            {showPDF && (
                <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center">
                    <div className="bg-white rounded-lg w-[90%] h-[90%]">
                        <div className="p-4 border-b flex justify-between items-center">
                            <h2 className="text-lg font-semibold">Link Performance Report</h2>
                            <Button variant="ghost" onClick={() => setShowPDF(false)}>
                                Close
                            </Button>
                        </div>
                        <div className="h-[calc(100%-4rem)]">
                            <PDFViewer width="100%" height="100%">
                                <LinkReportDocument links={links} />
                            </PDFViewer>
                        </div>
                    </div>
                </div>
            )}
        </>
    )
} 