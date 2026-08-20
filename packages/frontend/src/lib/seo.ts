/**
 * サイト全体で共有する SEO 関連の定数。
 * ページごとのメタタグ生成は composables/usePageMeta.ts が担う。
 */

export const SITE_URL = 'https://tiara-hakodate.com'
export const SITE_NAME = 'Tiara'
export const DEFAULT_TITLE = 'Tiara（函館 ニュークラブ ティアラ）'
export const DEFAULT_DESCRIPTION =
  '函館にあるニュークラブ「Tiara（ティアラ）」の公式ウェブサイトです。店舗情報やスタッフ、出勤スケジュール、料金システムをご覧いただけます。'
export const OG_IMAGE_URL = `${SITE_URL}/og-image.png`

/**
 * JSON-LD を <script> に埋め込む際のエスケープ。
 * 文字列中に "</script>" が含まれてもタグが閉じないよう "<" を Unicode エスケープする。
 */
export function toJsonLd(data: Record<string, unknown>): string {
  return JSON.stringify(data).replace(/</g, '\\u003c')
}
