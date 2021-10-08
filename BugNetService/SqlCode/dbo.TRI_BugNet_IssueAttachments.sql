
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		S-Setkov
-- Create date: 12.01.2021
-- Alter date:	08.10.2021 add [ProjectName] column
--				22.04.2021
-- Description:	Iserv sunc attachment
-- =============================================
CREATE OR ALTER TRIGGER dbo.TRI_BugNet_IssueAttachments ON dbo.BugNet_IssueAttachments
   AFTER INSERT
AS 
BEGIN
	SET NOCOUNT ON;

	INSERT	dbo.Iserv_MessageQueue(IssueId, TfsId, [User], [Operation], [Message], [AttachmentId], [FileName], [ProjectName])
	SELECT	i.IssueId, 
			i.TfsId, 
			ISNULL(up.DisplayName, u.UserName) AS [User],
			'add attachment' AS [Operation],
			'<a href="http://support.it-serv.ru/bugnet/Issues/UserControls/DownloadAttachment.axd?id=' + CAST(a.IssueAttachmentId AS VARCHAR(10)) + '">' 
			+ ISNULL(STUFF(a.FileName, PATINDEX('%[0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F]-[0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F]-[0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F]-[0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F]-[0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F]%', a.FileName), 37, ''), a.FileName)
			+ '</a>' AS [Message],
			a.IssueAttachmentId AS [AttachmentId],
			ISNULL(STUFF(a.FileName, PATINDEX('%[0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F]-[0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F]-[0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F]-[0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F]-[0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F]%', a.FileName), 37, ''), a.FileName) AS [FileName],
			p.ProjectName
	FROM INSERTED AS INS
		INNER JOIN dbo.BugNet_IssueAttachments a
			ON a.IssueAttachmentId = INS.IssueAttachmentId
		INNER JOIN dbo.BugNet_Issues i
			ON i.IssueId = a.IssueId
			AND i.TfsId <> 0
		INNER JOIN  dbo.BugNet_Projects p
			ON p.ProjectId = i.ProjectId
		INNER JOIN dbo.Users u
			ON u.UserId = a.UserId
		LEFT JOIN dbo.BugNet_UserProfiles up
			ON up.UserName = u.UserName

END
GO
