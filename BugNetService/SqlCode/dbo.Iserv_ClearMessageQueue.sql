SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- ===================================================================
-- Author:		S-Setkov
-- CREATE date: 25.12.2020
-- ALTER  date: 25.12.2020
-- Description:	Iserv clear message queue
-- ===================================================================
CREATE OR ALTER PROCEDURE dbo.Iserv_ClearMessageQueue
AS
BEGIN
	SET NOCOUNT, XACT_ABORT, ANSI_PADDING, ANSI_WARNINGS, ARITHABORT, CONCAT_NULL_YIELDS_NULL ON
	SET NUMERIC_ROUNDABORT, CURSOR_CLOSE_ON_COMMIT OFF

	DELETE	q
	FROM	dbo.Iserv_MessageQueue q
	WHERE	q.DateSync < DATEADD(WEEK, - 1, GETDATE())

END