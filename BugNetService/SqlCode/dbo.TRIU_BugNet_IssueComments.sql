
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		S-Setkov
-- Create date: 04.01.2021
-- Alter date:	29.09.2021 test update "Comment" column
--				12.01.2021
-- Description:	Iserv sunc comment
-- =============================================
CREATE OR ALTER TRIGGER dbo.TRIU_BugNet_IssueComments ON dbo.BugNet_IssueComments
   AFTER INSERT, UPDATE
AS 
BEGIN
	SET NOCOUNT ON;

	DECLARE @operation NVARCHAR(50)
	IF EXISTS(SELECT TOP 1 1 FROM DELETED)
		SET @operation = 'edit comment'
	ELSE
		SET @operation = 'add comment'

	IF UPDATE(Comment)
	BEGIN
		INSERT	dbo.Iserv_MessageQueue(IssueId, TfsId, [User], [Operation], [Message])
		SELECT	i.IssueId, 
				i.TfsId, 
				ISNULL(up.DisplayName, u.UserName) AS [User],
				@operation AS [Operation],
				c.Comment AS [Message]
		FROM INSERTED AS INS
			INNER JOIN dbo.BugNet_IssueComments c
				ON c.IssueCommentId = INS.IssueCommentID
			INNER JOIN dbo.BugNet_Issues i
				ON i.IssueId = c.IssueId
				AND i.TfsId <> 0
			INNER JOIN dbo.Users u
				ON u.UserId = c.UserId
			LEFT JOIN dbo.BugNet_UserProfiles up
				ON up.UserName = u.UserName
	END

END
GO
