import React from 'react'
import Box from '@material-ui/core/Box'
import { makeStyles, createStyles, Theme } from '@material-ui/core/styles'

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        container: {
            ...theme.typography.h4,
            width: 200,
            textAlign: 'center',
            color: theme.palette.primary.contrastText,
            shadow: theme.palette.warning.main,
        },
        suffix: {
            textShadow: `3px 0 ${theme.palette.primary.main}, -3px 0px ${theme.palette.primary.main}`,
            '&::after': {
                content: '""',
                width: '80%',
                marginLeft: 'auto',
                marginRight: 'auto',
                borderBottom: `3px solid ${theme.palette.warning.light}`,
                display: 'block',
                marginTop: -5,
            },
        },
    })
)

export default function Logo({}) {
    const classes = useStyles()
    return (
        <Box className={classes.container}>
            <span className={classes.suffix}>ov-chipkaart</span>
            dashboard
        </Box>
    )
}
