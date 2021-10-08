SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE dbo.Iserv_MessageQueue
(
	[Id] INT IDENTITY(1,1) NOT NULL,
	[Link] AS ([Id]) PERSISTED, -- deprecated
	[Date] SMALLDATETIME CONSTRAINT [DF_Iserv_MessageQueue_Date] DEFAULT (GETDATE()) NOT NULL,
	[IssueId] INT NOT NULL,
	[TfsId] INT NOT NULL,
	[User] NVARCHAR(100),
	[Operation] NVARCHAR(50),
	[Message] NVARCHAR(max),
	[DateSync] SMALLDATETIME,
	[IssueUrl] AS 'http://support.it-serv.ru/bugnet/Issues/IssueDetail.aspx?id=' + cast([IssueId] as NVARCHAR(250)) PERSISTED,
	[TfsUrl] AS 'http://tfs2017.compulink.local:8080/tfs/DefaultCollection/IServ/_workitems?id=' + cast([TfsId] as NVARCHAR(250)) PERSISTED,
	[AttachmentId] INT,
	[FileName] NVARCHAR(250),
	[ContentType] NVARCHAR(250),  -- deprecated
	[FileUrl] AS 'http://support.it-serv.ru/bugnet/Issues/UserControls/DownloadAttachment.axd?id=' + cast([AttachmentId] as NVARCHAR(250)) PERSISTED,
	[ProjectName] NVARCHAR(50),
	CONSTRAINT [PK_Iserv_MessageQueue] PRIMARY KEY CLUSTERED ([Id] ASC) ON [PRIMARY]
) ON [PRIMARY] 
GO