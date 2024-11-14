SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- ===================================================================
-- Author:		S-Setkov
-- CREATE date: 25.12.2020
-- ALTER  date: 14.11.2024
-- Description:	Iserv clear message queue
-- ===================================================================
CREATE OR ALTER PROCEDURE dbo.Iserv_ClearMessageQueue
AS
BEGIN
	SET NOCOUNT, XACT_ABORT, ANSI_PADDING, ANSI_WARNINGS, ARITHABORT, CONCAT_NULL_YIELDS_NULL ON
	SET NUMERIC_ROUNDABORT, CURSOR_CLOSE_ON_COMMIT OFF

--	удалить записи старше месяца
	DELETE	q
	FROM	dbo.Iserv_MessageQueue q
	WHERE	q.DateSync < DATEADD(MONTH, - 1, GETDATE())

--	вырезать вложенные картинки из сообщения
	WHILE EXISTS(
			SELECT	TOP 1 q.Id
			FROM	dbo.Iserv_MessageQueue q
			WHERE	q.DateSync IS NOT NULL
				AND CHARINDEX('<img', Message) > 0
			)
	BEGIN
		UPDATE	q
			SET Message = CONCAT(LEFT(Message, i.start_pos-1), '<p> ... тут была картинка ... </p>', RIGHT(Message, LEN(Message)-i.end_pos)) 
		FROM	dbo.Iserv_MessageQueue q
				OUTER APPLY (
					SELECT	CHARINDEX('<img', Message) AS start_pos,
							CHARINDEX('/>', Message, CHARINDEX('<img', Message)) + 1 AS end_pos
							) AS i
		WHERE	q.DateSync IS NOT NULL
			AND	i.start_pos > 0
	END

END