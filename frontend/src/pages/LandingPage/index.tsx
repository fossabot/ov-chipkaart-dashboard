import React from 'react'
import './index.scss'
import Grid from '@material-ui/core/Grid'
import { makeStyles, createStyles, Theme } from '@material-ui/core/styles'
import Box from '@material-ui/core/Box'

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        container: {
            height: '100vh',
            width: '100vw',
            overflow: 'hidden',
        },

        info: {
            height: '100vh',
            display: 'flex',
            alignItems: 'center',
            padding: theme.spacing(2),
            backgroundColor: theme.palette.primary.dark,
        },

        auth: {
            height: '100vh',
            display: 'flex',
            alignItems: 'center',
            padding: theme.spacing(3),
            backgroundColor: theme.palette.primary.contrastText,
        },
    })
)

export default function LandingPage() {
    const classes = useStyles()
    return (
        <Grid container className={classes.container}>
            <Grid item xs={8}>
                <Box className={classes.info}>First</Box>
            </Grid>
            <Grid item xs={4}>
                <Box className={classes.auth}>Second</Box>
            </Grid>
        </Grid>
    )
}
